package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WorkRequest struct {
	URL     url.URL
	Execute func(workerID, numberOfWorkers int) WorkResult
}

var numberOfWorkers int

func executeWork(workerID, numberOfWorkers int, targetURL url.URL, newURLs chan url.URL) WorkResult {

	// read the URL
	response, err := readURL(targetURL)
	if err != nil {
		return WorkResult{err: err}
	}

	if response.IsHTML() {

		// get dependent links
		links, err := getDependentRequests(targetURL, bytes.NewReader(response.Body()))
		if err != nil {
			return WorkResult{err: err}
		}

		for _, link := range links {
			go func(url url.URL) {
				newURLs <- url
			}(link)
		}

	}

	workResult := WorkResult{
		url: targetURL,

		workerID:        workerID,
		numberOfWorkers: numberOfWorkers,

		responseSize: response.Size(),
		statusCode:   response.StatusCode(),
		startTime:    response.StartTime(),
		endTime:      response.EndTime(),
		contentType:  response.ContentType(),
	}

	return workResult
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
			href = strings.TrimSuffix(base, "/") + "/" + href
		}

		// cut off hashes
		if strings.Contains(href, "#") {
			hashPosition := strings.Index(href, "#")
			href = href[0 : hashPosition-1]
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
