package main

import (
	"github.com/pkg/errors"
	"log"
	"net/url"
	"os"
	"time"
)

type CrawlOptions struct {
	NumberOfConcurrentRequests int
	Timeout                    time.Duration
	UserAgent                  string
	LogFile                    string
}

func crawl(xmlSitemapURL url.URL, options CrawlOptions, stop chan bool) error {

	// read the XML sitemap as a initial source for URLs
	urlsFromXMLSitemap, err := getURLs(xmlSitemapURL, "gargantua bot")
	if err != nil {
		return err
	}

	// the URL queue
	urls := make(chan crawlerUrl, len(urlsFromXMLSitemap))

	// fill the URL queue with the URLs from the XML sitemap
	for _, xmlSitemapURLEntry := range urlsFromXMLSitemap {
		urls <- xmlSitemapURLEntry
	}

	results := make(chan WorkResult)

	// send new urls to the work queue
	workers := make(chan int, options.NumberOfConcurrentRequests)
	for workerID := 1; workerID <= options.NumberOfConcurrentRequests; workerID++ {
		workers <- workerID
	}

	allURLsHaveBeenVisited := make(chan bool)
	go func() {
		var visitedURLs = make(map[string]crawlerUrl)
		for {
			select {
			case <-stop:
				allURLsHaveBeenVisited <- true
				return

			case targetURL := <-urls:
				// skip URLs we have already seen
				_, alreadyVisited := visitedURLs[targetURL.String()]

				if alreadyVisited {
					continue
				}

				// mark the URL as visited
				visitedURLs[targetURL.String()] = targetURL

				debugf("Sending URL to work queue: %s", targetURL.String())

				go func() {
					workerID := <-workers
					debugf("Using worker %d for URL %q", workerID, targetURL.String())
					results <- executeWork(workerID, cap(workers), targetURL, options.UserAgent, urls)
					debugf("Worker %d finished processing URL %q", workerID, targetURL.String())
					workers <- workerID
				}()

			case <-time.After(time.Second * 1):

				if len(workers) == cap(workers) && len(urls) == 0 {
					allURLsHaveBeenVisited <- true
					return
				}

			}
		}
	}()

	var logger *log.Logger
	if options.LogFile != "" {
		file, err := os.OpenFile(options.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return errors.Wrapf(err, "failed to open log file %q for writing", options.LogFile)
		}

		defer file.Close()
		logger = log.New(file, "", log.Ldate|log.Ltime)
	}

	// update the statistics with the results
	allStatisticsHaveBeenUpdated := make(chan bool)
	go func() {
		for {
			select {
			case <-allURLsHaveBeenVisited:
				allStatisticsHaveBeenUpdated <- true
				return

			case result := <-results:
				receivedUrl := result.URL()
				debugf("Received results for URL %q", receivedUrl.String())
				updateStatistics(result)

				if logger != nil {
					logResult(logger, result)
				}
			}
		}
	}()

	<-allStatisticsHaveBeenUpdated

	return nil
}
