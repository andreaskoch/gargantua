package main

import (
	"encoding/xml"
	"strings"
)

func getXMLSitemap(xmlSitemapURL string) (XMLSitemap, error) {

	response, readErr := readURL(xmlSitemapURL)
	if readErr != nil {
		return XMLSitemap{}, readErr
	}

	if !strings.Contains(string(response.Body()), "</urlset>") {
		return XMLSitemap{}, XmlSitemapError{"Invalid content"}
	}

	var urlSet XMLSitemap
	unmarshalError := xml.Unmarshal(response.Body(), &urlSet)
	if unmarshalError != nil {
		return XMLSitemap{}, unmarshalError
	}

	return urlSet, nil
}

type XMLSitemap struct {
	URLs []URL `xml:"url"`
}

type URL struct {
	Location string `xml:"loc"`
}

type XmlSitemapError struct {
	message string
}

func (sitemapIndexError XmlSitemapError) Error() string {
	return sitemapIndexError.message
}

func isInvalidXMLSitemapContent(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "Invalid content"
}
