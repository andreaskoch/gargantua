package main

import (
	"bytes"
	"net/url"
)

type WorkRequest struct {
	URL     url.URL
	Execute func(workerID, numberOfWorkers int) WorkResult
}

func executeWork(workerID, numberOfWorkers int, targetURL url.URL, userAgent string, newURLs chan url.URL) WorkResult {

	// read the URL
	response, err := readURL(targetURL, userAgent)
	if err != nil {
		return WorkResult{err: err}
	}

	if response.IsHTML() {

		// get dependent links
		links, err := getDependentRequests(targetURL, bytes.NewReader(response.Body()))
		if err != nil {
			return WorkResult{err: err}
		}

		for _, link := range links {
			go func(url url.URL) {
				newURLs <- url
			}(link)
		}

	}

	workResult := WorkResult{
		url: targetURL,

		workerID:        workerID,
		numberOfWorkers: numberOfWorkers,

		responseSize: response.Size(),
		statusCode:   response.StatusCode(),
		startTime:    response.StartTime(),
		endTime:      response.EndTime(),
		contentType:  response.ContentType(),
	}

	return workResult
}
