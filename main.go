package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const applicationName = "gargantua"
const applicationVersion = "v0.1.0-alpha"

var (
	app = kingpin.New(applicationName, fmt.Sprintf(`「 %s 」%s crawls all URLs of your website - starting with the links in your sitemap.xml

    🌈 https://github.com/andreaskoch/gargantua
`, applicationName, applicationVersion))

	// global
	verbose = app.Flag("verbose", "Disable the UI and enable debug mode").Envar("GARGANTUA_VERBOSE").Short('v').Default("false").Bool()
	timeout = app.Flag("timeout", "The HTTP timeout in seconds used by the crawler").Envar("GARGANTUA_TIMEOUT").Short('t').Default("60").Int()

	// crawl
	crawlCommand    = app.Command("crawl", "Crawls a given websites' XML sitemap")
	crawlWebsiteURL = crawlCommand.Flag("url", "The URL to a websites' XML sitemap (e.g. https://www.sitemaps.org/sitemap.xml)").Required().Envar("GARGANTUA_URL").Short('u').String()
	crawlWorkers    = crawlCommand.Flag("workers", "The number of concurrent workers that crawl the site at the same time").Required().Envar("GARGANTUA_WORKERS").Short('w').Int()
)

func init() {
	app.Version(applicationVersion)
	app.Author("Andreas Koch <andy@ak7.io>")
}

func main() {
	handleCommandlineArgument(os.Args[1:])
}

func handleCommandlineArgument(arguments []string) {

	switch kingpin.MustParse(app.Parse(arguments)) {

	case crawlCommand.FullCommand():
		websiteURL, parseError := url.Parse(*crawlWebsiteURL)
		if parseError != nil {
			fmt.Fprintf(os.Stderr, "%s", parseError.Error())
			os.Exit(1)
		}

		err := startCrawling(*websiteURL, *crawlWorkers, *timeout, *verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

}

func startCrawling(targetURL url.URL, concurrentRequests, timeoutInSeconds int, debugModeIsEnabled bool) error {
	stopTheCrawler := make(chan bool)
	crawlResult := make(chan error)

	go func() {
		result := crawl(targetURL, CrawlOptions{
			NumberOfConcurrentRequests: int(concurrentRequests),
			Timeout:                    time.Second * time.Duration(timeoutInSeconds),
		}, stopTheCrawler)

		crawlResult <- result
	}()

	var uiWaitGroup = &sync.WaitGroup{}
	if debugModeIsEnabled {
		debugf = consoleDebug
	} else {
		debugf = dashboardDebug

		uiWaitGroup.Add(1)
		go func() {
			dashboard(time.Now(), stopTheCrawler)
			uiWaitGroup.Done()
		}()
	}

	uiWaitGroup.Wait()

	err := <-crawlResult
	if err != nil {
		return err
	}

	return nil
}
