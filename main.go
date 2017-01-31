package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	targetURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %q: %s", os.Args[1], err)
		os.Exit(1)
	}

	crawl(*targetURL, CrawlOptions{})
}
