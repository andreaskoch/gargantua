package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
)

var applicationName = "sitemap-crawler"

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

	crawl(*targetURL, CrawlOptions{int(numberOfParallelRequests)})
}

func usage(writer io.Writer) {
	fmt.Fprintf(writer, "«%s» crawls your website starting the URLs in your XML sitemap\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "Usage:\n\n")
	fmt.Fprintf(writer, "  %s <number-of-parallel-requests> <sitemap-url>\n\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "Example:\n\n")
	fmt.Fprintf(writer, "  %s 20 http://example.com/sitemap.org\n\n", applicationName)
}
