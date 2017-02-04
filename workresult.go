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
