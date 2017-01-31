package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var urlLock = sync.RWMutex{}
var visitedURLs = make(map[string]url.URL)

type URLCrawlRequest struct {
	ParentURL url.URL
	TargetURL url.URL
}

func crawl(xmlSitemapURL url.URL, options CrawlOptions) error {

	// read the XML sitemap as a initial source for URLs
	urlsFromXMLSitemap, err := getURLs(xmlSitemapURL)
	if err != nil {
		return err
	}

	// the URL queue
	crawlRequestQueue := make(chan URLCrawlRequest, len(urlsFromXMLSitemap))

	// fill the URL queue with the URLs from the XML sitemap
	for _, xmlSitemapURLEntry := range urlsFromXMLSitemap {
		crawlRequestQueue <- URLCrawlRequest{
			ParentURL: xmlSitemapURL,
			TargetURL: xmlSitemapURLEntry,
		}
	}

	// crawl all URLs in the queue
	results := StartDispatcher(50)
	go func() {
		for urlCrawlRequest := range crawlRequestQueue {

			// skip URLs we have already seen
			urlLock.RLock()
			_, alreadyVisited := visitedURLs[urlCrawlRequest.TargetURL.String()]
			urlLock.RUnlock()

			if alreadyVisited {
				continue
			}

			// mark the URL as visited
			urlLock.Lock()
			visitedURLs[urlCrawlRequest.TargetURL.String()] = urlCrawlRequest.TargetURL
			urlLock.Unlock()

			go func(urlToCrawl URLCrawlRequest) {
				WorkQueue <- createWorkRequest(urlToCrawl, crawlRequestQueue)
			}(urlCrawlRequest)

		}
	}()

	// present the results
	doneCounter := 0
	for result := range results {

		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "%s\n", result.Error)
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", result.Message)
		}

		doneCounter++
	}

	return nil
}

func createWorkRequest(urlToCrawl URLCrawlRequest, newUrls chan URLCrawlRequest) WorkRequest {

	return WorkRequest{
		URL: urlToCrawl.TargetURL,
		Execute: func() WorkResult {

			// read the URL
			response, err := readURL(urlToCrawl.TargetURL)
			if err != nil {
				return WorkResult{
					Error: err,
				}
			}

			if response.IsHTML() {

				// get dependent links
				links, err := getDependentRequests(urlToCrawl.TargetURL, bytes.NewReader(response.Body()))
				if err != nil {
					return WorkResult{
						Error: err,
					}
				}

				for _, link := range links {
					newUrls <- URLCrawlRequest{
						ParentURL: urlToCrawl.TargetURL,
						TargetURL: link,
					}
				}

			}

			return WorkResult{
				Message: fmt.Sprintf("%03d %9s %15s  %s  %s", response.StatusCode(), fmt.Sprintf("%d", response.Size()), fmt.Sprintf("%s", response.Duration()), urlToCrawl.ParentURL.String(), urlToCrawl.TargetURL.String()),
			}
		}}

}

func getDependentRequests(baseURL url.URL, input io.Reader) ([]url.URL, error) {

	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		return nil, err
	}

	var urls []url.URL

	// base url
	base, _ := doc.Find("base[href]").Attr("href")
	if base == "" {
		base = baseURL.Scheme + "://" + baseURL.Host
	} else if strings.HasPrefix(base, "/") {
		base = baseURL.Scheme + "://" + baseURL.Host + base
	}

	// get all links
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, _ := s.Attr("href")

		if strings.HasPrefix(href, "/") {
			href = baseURL.Scheme + "://" + baseURL.Host + href
		} else if !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") {
			href = strings.TrimSuffix(base, "/") + href
		}

		hrefURL, err := url.Parse(href)
		if err != nil {
			return
		}

		// ignore external links
		if hrefURL.Host != baseURL.Host {
			return
		}

		urls = append(urls, *hrefURL)
	})

	return urls, nil
}

func getURLs(xmlSitemapURL url.URL) ([]url.URL, error) {

	var urls []url.URL

	urlsFromIndex, indexError := getURLsFromSitemapIndex(xmlSitemapURL)
	if indexError == nil {
		urls = urlsFromIndex
	}

	urlsFromSitemap, sitemapError := getURLsFromSitemap(xmlSitemapURL)
	if sitemapError == nil {
		urls = append(urls, urlsFromSitemap...)
	}

	if isInvalidSitemapIndexContent(indexError) && isInvalidXMLSitemapContent(sitemapError) {
		return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", xmlSitemapURL.String())
	}

	return urls, nil

}

func getURLsFromSitemap(xmlSitemapURL url.URL) ([]url.URL, error) {

	var urls []url.URL

	sitemap, xmlSitemapError := getXMLSitemap(xmlSitemapURL)
	if xmlSitemapError != nil {
		return nil, xmlSitemapError
	}

	for _, urlEntry := range sitemap.URLs {

		parsedURL, parseError := url.Parse(urlEntry.Location)
		if parseError != nil {
			return nil, parseError
		}

		urls = append(urls, *parsedURL)
	}

	return urls, nil
}

func getURLsFromSitemapIndex(xmlSitemapURL url.URL) ([]url.URL, error) {

	var urls []url.URL

	sitemapIndex, sitemapIndexError := getSitemapIndex(xmlSitemapURL)
	if sitemapIndexError != nil {
		return nil, sitemapIndexError
	}

	for _, sitemap := range sitemapIndex.Sitemaps {

		locationURL, err := url.Parse(sitemap.Location)
		if err != nil {
			return nil, err
		}

		sitemapUrls, err := getURLsFromSitemap(*locationURL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, sitemapUrls...)
	}

	return urls, nil

}

type CrawlOptions struct {
	Hosts []net.IP
}
