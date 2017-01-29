package main

import (
	"fmt"
	"net"
	"net/url"
)

func crawl(xmlSitemapURL string, options CrawlOptions) error {

	urls, err := getURLs(xmlSitemapURL)
	if err != nil {
		return err
	}

	results := StartDispatcher(80, len(urls))

	for index, url := range urls {

		// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
		work := WorkRequest{Name: fmt.Sprintf("%000d %s", index+1, url.String()), Execute: createWorkFunction(index, url)}

		// Push the work onto the queue.
		go func() {
			WorkQueue <- work
		}()
	}

	for result := range results {
		fmt.Println(result.Message)
	}

	return nil
}

func createWorkFunction(index int, url url.URL) func() WorkResult {
	return func() WorkResult {

		content, err := readURL(url.String())
		if err != nil {
			return WorkResult{fmt.Sprintf("Error: %s", err)}
		}

		return WorkResult{fmt.Sprintf("%000d %s: %d", index+1, url.String(), len(content))}
	}
}

func getURLs(xmlSitemapURL string) ([]url.URL, error) {

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
		return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", xmlSitemapURL)
	}

	return urls, nil

}

func getURLsFromSitemap(xmlSitemapURL string) ([]url.URL, error) {

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

func getURLsFromSitemapIndex(xmlSitemapURL string) ([]url.URL, error) {

	var urls []url.URL

	sitemapIndex, sitemapIndexError := getSitemapIndex(xmlSitemapURL)
	if sitemapIndexError != nil {
		return nil, sitemapIndexError
	}

	for _, sitemap := range sitemapIndex.Sitemaps {

		sitemapUrls, err := getURLsFromSitemap(sitemap.Location)
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
