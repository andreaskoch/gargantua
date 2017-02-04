package main

import (
	"net/url"
	"testing"
)

func Test_getShortenedURL(t *testing.T) {

	url, _ := url.Parse("https://example.com/part1/part2/part3/file.html")
	result := getShortenedURL(*url, 15)

	if result != "part3/file.html" {
		t.Fail()
	}
}
