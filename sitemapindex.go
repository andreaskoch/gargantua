package main

import (
	"encoding/xml"
	"net/url"
	"strings"
)

func getSitemapIndex(xmlSitemapURL url.URL, userAgent string) (SitemapIndex, error) {
	response, readErr := readURL(xmlSitemapURL, userAgent)
	if readErr != nil {
		return SitemapIndex{}, readErr
	}

	if !strings.Contains(string(response.Body()), "</sitemapindex>") {
		return SitemapIndex{}, SitemapIndexError{"Invalid content"}
	}

	var sitemapIndex SitemapIndex
	unmarshalError := xml.Unmarshal(response.Body(), &sitemapIndex)
	if unmarshalError != nil {
		return SitemapIndex{}, unmarshalError
	}

	return sitemapIndex, nil
}

type SitemapIndex struct {
	Sitemaps []URL `xml:"sitemap"`
}

type Sitemap struct {
	Location string `xml:"loc"`
}

type SitemapIndexError struct {
	message string
}

func (sitemapIndexError SitemapIndexError) Error() string {
	return sitemapIndexError.message
}

func isInvalidSitemapIndexContent(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "Invalid content"
}
