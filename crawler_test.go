package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getTestSitemapXml(domain string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
  <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1">

  <url>
    <loc>http://%s/development</loc>
    <lastmod>2015-01-04</lastmod>
    <changefreq>never</changefreq>
    <priority>1.0</priority>
  </url>

  </urlset>`, domain)
}

func Test_crawl_validSitemap(t *testing.T) {

	testContentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "Some Content")

	}))

	contentServerURL, _ := url.Parse(testContentServer.URL)

	defer testContentServer.Close()

	testSitemapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, getTestSitemapXml(contentServerURL.Host))

	}))

	defer testSitemapServer.Close()

	crawl(testSitemapServer.URL, CrawlOptions{})
}
