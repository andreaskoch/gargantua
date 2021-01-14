package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type WorkResult struct {
	err error

	parentURL url.URL
	url       url.URL

	numberOfWorkers int
	workerID        int

	responseSize int
	body         []byte
	header       map[string][]string

	statusCode  int
	startTime   time.Time
	endTime     time.Time
	contentType string
}

func (workResult WorkResult) String() string {

	headers := []string{}
	for name, value := range workResult.header {
		headers = append(headers, fmt.Sprintf("'%s:%s'", name, value))
	}

	return fmt.Sprintf("#%03d: %03d %9s %15s %20s %20s %s",
		workResult.workerID,
		workResult.statusCode,
		fmt.Sprintf("%d", workResult.responseSize),
		fmt.Sprintf("%fms", workResult.ResponseTime().Seconds()*1000),
		workResult.url.String(),
		workResult.parentURL.String(),
		strings.Join(headers, "|"),
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

func (workResult WorkResult) Body() []byte {
	return workResult.body
}

func (workResult WorkResult) Header() map[string][]string {
	return workResult.header
}
