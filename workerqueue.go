package main

import (
	"fmt"
	"net/url"
	"time"
)

type WorkResult struct {
	err error

	parentURL url.URL
	url       url.URL

	numberOfWorkers int
	workerID        int

	responseSize int
	statusCode   int
	startTime    time.Time
	endTime      time.Time
	contentType  string
}

func (workResult WorkResult) String() string {
	return fmt.Sprintf("#%03d: %03d %9s %15s %20s",
		workResult.workerID,
		workResult.statusCode,
		fmt.Sprintf("%d", workResult.responseSize),
		fmt.Sprintf("%f ms", workResult.ResponseTime().Seconds()*1000),
		workResult.url.String(),
	)
}

func (workResult WorkResult) Error() error {
	return workResult.err
}

func (workResult WorkResult) ParentURL() url.URL {
	return workResult.parentURL
}

func (workResult WorkResult) URL() url.URL {
	return workResult.url
}

func (workResult WorkResult) Size() int {
	return workResult.responseSize
}

func (workResult WorkResult) StatusCode() int {
	return workResult.statusCode
}

func (workResult WorkResult) StartTime() time.Time {
	return workResult.startTime
}

func (workResult WorkResult) EndTime() time.Time {
	return workResult.endTime
}

func (workResult WorkResult) ResponseTime() time.Duration {
	return workResult.endTime.Sub(workResult.startTime)
}

func (workResult WorkResult) ContentType() string {
	return workResult.contentType
}

func (workResult WorkResult) WorkerID() int {
	return workResult.workerID
}

func (workResult WorkResult) NumberOfWorkers() int {
	return workResult.numberOfWorkers
}

type WorkRequest struct {
	URL     url.URL
	Execute func(workerID, numberOfWorkers int) WorkResult
}

var WorkerQueue chan chan WorkRequest

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

func StartDispatcher(numberOfWorkers int) chan WorkResult {

	results := make(chan WorkResult, 10)

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, numberOfWorkers)

	// Now, create all of our workers.
	for i := 0; i < numberOfWorkers; i++ {
		worker := NewWorker(i+1, WorkerQueue, results)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()

	return results
}
