package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_getXMLSitemap_InvalidContent_ErrorIsReturned(t *testing.T) {
	testSitemapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "li da da")

	}))

	defer testSitemapServer.Close()

	testServerURL, _ := url.Parse(testSitemapServer.URL)
	_, err := getXMLSitemap(*testServerURL)

	if err == nil {
		t.Fail()
		t.Logf("getXMLSitemap(%q) should have returned an error", testSitemapServer.URL)
	}

}

func Test_getXMLSitemap_SitemapExists_SitemapIsReturned(t *testing.T) {
	testSitemapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, getTestSitemapXml("example.com"))

	}))

	defer testSitemapServer.Close()

	testServerURL, _ := url.Parse(testSitemapServer.URL)
	sitemap, err := getXMLSitemap(*testServerURL)

	if err != nil {
		t.Fail()
		t.Logf("getXMLSitemap(%q) returned an error: %s", testSitemapServer.URL, err)
	}

	if len(sitemap.URLs) == 0 {
		t.Fail()
		t.Logf("getXMLSitemap(%q) returned 0 URLs", testSitemapServer.URL)
	}

}
