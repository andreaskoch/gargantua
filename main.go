package main

import "os"

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	crawl(os.Args[1], CrawlOptions{})
}
