package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_getSitemapIndex_NoIndexGiven_ErrorIsReturned(t *testing.T) {
	sitemapIndexContent := `la di da`

	testSitemapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, sitemapIndexContent)

	}))

	defer testSitemapServer.Close()

	testServerURL, _ := url.Parse(testSitemapServer.URL)
	_, err := getSitemapIndex(*testServerURL)

	if err == nil {
		t.Fail()
		t.Logf("getSitemapIndex(%q) should have returned an error", testSitemapServer.URL)
	}

}

func Test_getSitemapIndex_IndexExists_IndexIsNotEmpty(t *testing.T) {
	sitemapIndexContent := `<?xml version="1.0" encoding="UTF-8"?>
  <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <sitemap>
      <loc>https://example.com/sitemap.xml</loc>
      <lastmod>2016-10-06T11:20:05+00:00</lastmod>
    </sitemap>
  </sitemapindex>`

	testSitemapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, sitemapIndexContent)

	}))

	defer testSitemapServer.Close()

	testServerURL, _ := url.Parse(testSitemapServer.URL)
	sitemapIndex, err := getSitemapIndex(*testServerURL)

	if err != nil {
		t.Fail()
		t.Logf("getSitemapIndex(%q) returned %s", testSitemapServer.URL, err)
	}

	if len(sitemapIndex.Sitemaps) == 0 {
		t.Fail()
		t.Logf("getSitemapIndex(%q) does not contain any sitemaps", testSitemapServer.URL)
	}

}
