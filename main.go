package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"
)

const applicationName = "gargantua"
const applicationVersion = "v0.1.0-alpha"

const minimumNumberOfConcurrentRequests = 1
const maxiumumNumberOfConcurrentRequests = 1000

func main() {
	if len(os.Args) < 3 {
		usage(os.Stderr)
		os.Exit(1)
	}

	numberOfConcurrentRequests, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %q as the number of concurrent requests: %s", os.Args[1], err)
		os.Exit(1)
	}

	if numberOfConcurrentRequests < minimumNumberOfConcurrentRequests || numberOfConcurrentRequests > maxiumumNumberOfConcurrentRequests {
		fmt.Fprintf(os.Stderr, "The number of concurrent requests must be between %d amd %d", minimumNumberOfConcurrentRequests, maxiumumNumberOfConcurrentRequests)
		os.Exit(1)
	}

	targetURL, err := url.Parse(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %q as the XML sitemap URL: %s", os.Args[2], err)
		os.Exit(1)
	}

	stopTheCrawler := make(chan bool)
	crawlResult := make(chan error)

	go func() {
		result := crawl(*targetURL, CrawlOptions{
			NumberOfConcurrentRequests: int(numberOfConcurrentRequests),
			Timeout:                    time.Second * 10,
		}, stopTheCrawler)

		crawlResult <- result
	}()

	interactiveUI := false

	debugf = consoleDebug
	if interactiveUI {
		debugf = dashboardDebug
		go dashboard(time.Now(), stopTheCrawler)
	}

	<-crawlResult

}

func usage(writer io.Writer) {
	fmt.Fprintf(writer, "「 %s 」crawls all URLs of your website - starting with the links in your sitemap.xml\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "%s\n\n", applicationVersion)
	fmt.Fprintf(writer, "Usage:\n\n")
	fmt.Fprintf(writer, "  %s <number-of-concurrent-requests> <sitemap-url>\n\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "Example:\n\n")
	fmt.Fprintf(writer, "  %s 20 http://example.com/sitemap.org\n\n", applicationName)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "🌈 https://github.com/andreaskoch/gargantua\n")
}

var errorMessages []string

var debugf func(format string, a ...interface{})

func init() {
	debugf = func(format string, a ...interface{}) {
		// default
	}
}

func consoleDebug(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func dashboardDebug(format string, a ...interface{}) {
	latestMesasges, err := getLatestLogMessages(errorMessages, 4)
	if err != nil {
		panic(err)
	}

	errorMessages = append(latestMesasges, fmt.Sprintf(format, a...))
}
