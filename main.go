package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
)

const applicationName = "gargantua"
const applicationVersion = "v0.1.0-alpha"

const minimumNumberOfParallelRequests = 1
const maxiumumNumberOfParallelRequests = 1000

func main() {
	if len(os.Args) < 3 {
		usage(os.Stderr)
		os.Exit(1)
	}

	numberOfParallelRequests, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %q as the number of parallel requests: %s", os.Args[1], err)
		os.Exit(1)
	}

	if numberOfParallelRequests < minimumNumberOfParallelRequests || numberOfParallelRequests > maxiumumNumberOfParallelRequests {
		fmt.Fprintf(os.Stderr, "The number of parallel requests must be between %d amd %d", minimumNumberOfParallelRequests, maxiumumNumberOfParallelRequests)
		os.Exit(1)
	}

	targetURL, err := url.Parse(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %q as the XML sitemap URL: %s", os.Args[2], err)
		os.Exit(1)
	}

	dashboard()

	done := make(chan bool)
	go func() {
		crawl(*targetURL, CrawlOptions{int(numberOfParallelRequests)})
		done <- true
	}()

	<-done

}

func usage(writer io.Writer) {
	fmt.Fprintf(writer, "「 %s 」crawls all URLs of your website - starting with the links in your sitemap.xml\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "%s\n\n", applicationVersion)
	fmt.Fprintf(writer, "Usage:\n\n")
	fmt.Fprintf(writer, "  %s <number-of-parallel-requests> <sitemap-url>\n\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "Example:\n\n")
	fmt.Fprintf(writer, "  %s 20 http://example.com/sitemap.org\n\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "🌈 https://github.com/andreaskoch/gargantua\n")
}
