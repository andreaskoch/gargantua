package main

import (
	"bytes"
	"net/url"
)

type WorkRequest struct {
	URL     url.URL
	Execute func(workerID, numberOfWorkers int) WorkResult
}

func executeWork(workerID, numberOfWorkers int, targetURL crawlerUrl, userAgent string, newURLs chan crawlerUrl) WorkResult {

	// read the URL
	response, err := readURL(targetURL.url, userAgent)
	if err != nil {
		return WorkResult{err: err}
	}

	if response.IsHTML() {

		// get dependent links
		links, err := getDependentRequests(targetURL.url, bytes.NewReader(response.Body()))
		if err != nil {
			return WorkResult{err: err}
		}

		for _, link := range links {
			go func(url crawlerUrl) {
				newURLs <- url
			}(link)
		}

	}

	workResult := WorkResult{
		parentURL: targetURL.parent,
		url:       targetURL.url,

		workerID:        workerID,
		numberOfWorkers: numberOfWorkers,

		responseSize: response.Size(),
		body:         response.Body(),
		header:       response.Header(),

		statusCode:  response.StatusCode(),
		startTime:   response.StartTime(),
		endTime:     response.EndTime(),
		contentType: response.ContentType(),
	}

	return workResult
}
