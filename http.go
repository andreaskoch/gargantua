package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Response struct {
	body        []byte
	statusCode  int
	startTime   time.Time
	endTime     time.Time
	contentType string
}

func (response *Response) Body() []byte {
	return response.body
}

func (response *Response) Size() int {
	return len(response.body)
}

func (response *Response) StatusCode() int {
	return response.statusCode
}

func (response *Response) StartTime() time.Time {
	return response.startTime
}

func (response *Response) EndTime() time.Time {
	return response.endTime
}

func (response *Response) ContentType() string {
	return response.contentType
}

func (response *Response) IsHTML() bool {
	return strings.HasPrefix(response.contentType, "text/html")
}

func readURL(url url.URL, userAgent string) (Response, error) {
	startTime := time.Now().UTC()

	req, requestErr := http.NewRequest("GET", url.String(), nil)
	if requestErr != nil {
		return Response{}, requestErr
	}

	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return Response{}, fetchErr
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return Response{}, readErr
	}

	endTime := time.Now().UTC()

	// content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	return Response{
		body:        body,
		statusCode:  resp.StatusCode,
		startTime:   startTime,
		endTime:     endTime,
		contentType: contentType,
	}, nil
}

func getURLs(xmlSitemapURL url.URL, userAgent string) ([]url.URL, error) {

	var urls []url.URL

	urlsFromIndex, indexError := getURLsFromSitemapIndex(xmlSitemapURL, userAgent)
	if indexError == nil {
		urls = urlsFromIndex
	}

	urlsFromSitemap, sitemapError := getURLsFromSitemap(xmlSitemapURL, userAgent)
	if sitemapError == nil {
		urls = append(urls, urlsFromSitemap...)
	}

	if isInvalidSitemapIndexContent(indexError) && isInvalidXMLSitemapContent(sitemapError) {
		return nil, fmt.Errorf("%q is neither a sitemap index nor a XML sitemap", xmlSitemapURL.String())
	}

	return urls, nil

}

func getURLsFromSitemap(xmlSitemapURL url.URL, userAgent string) ([]url.URL, error) {

	var urls []url.URL

	sitemap, xmlSitemapError := getXMLSitemap(xmlSitemapURL, userAgent)
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

func getURLsFromSitemapIndex(xmlSitemapURL url.URL, userAgent string) ([]url.URL, error) {

	var urls []url.URL

	sitemapIndex, sitemapIndexError := getSitemapIndex(xmlSitemapURL, userAgent)
	if sitemapIndexError != nil {
		return nil, sitemapIndexError
	}

	for _, sitemap := range sitemapIndex.Sitemaps {

		locationURL, err := url.Parse(sitemap.Location)
		if err != nil {
			return nil, err
		}

		sitemapUrls, err := getURLsFromSitemap(*locationURL, userAgent)
		if err != nil {
			return nil, err
		}

		urls = append(urls, sitemapUrls...)
	}

	return urls, nil

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
