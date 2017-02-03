package main

import (
	"sync"
	"time"
)

var stats Statistics

func init() {
	stats = Statistics{
		lock: sync.RWMutex{},
		numberOfRequestsByStatusCode:  make(map[int]int),
		numberOfRequestsByContentType: make(map[string]int),
	}
}

func updateStatistics(workResult WorkResult) {
	go stats.Add(workResult)
}

type SnapShot struct {

	// time
	timestamp           time.Time
	averageResponseTime time.Duration

	// counters
	totalNumberOfRequests        int
	numberOfSuccessfulRequests   int
	numberOfUnsuccessfulRequests int
	numberOfRequestsPerSecond    float64

	numberOfRequestsByStatusCode  map[int]int
	numberOfRequestsByContentType map[string]int

	// size
	totalSizeInBytes   int
	averageSizeInBytes int
}

type Statistics struct {
	lock sync.RWMutex

	rawResults []WorkResult
	snapShots  []SnapShot

	startTime time.Time
	endTime   time.Time

	totalResponseTime time.Duration

	numberOfRequests              int
	numberOfSuccessfulRequests    int
	numberOfUnsuccessfulRequests  int
	numberOfRequestsByStatusCode  map[int]int
	numberOfRequestsByContentType map[string]int

	totalSizeInBytes int
}

func (statistics *Statistics) Add(workResult WorkResult) SnapShot {

	// update the raw results
	statistics.lock.Lock()
	defer statistics.lock.Unlock()
	statistics.rawResults = append(statistics.rawResults, workResult)

	// initialize start and end time
	if statistics.numberOfRequests == 0 {
		statistics.startTime = workResult.StartTime()
		statistics.endTime = workResult.EndTime()
	}

	// start time
	if workResult.StartTime().Before(statistics.startTime) {
		statistics.startTime = workResult.StartTime()
	}

	// end time
	if workResult.EndTime().After(statistics.endTime) {
		statistics.endTime = workResult.EndTime()
	}

	// update the total number of requests
	statistics.numberOfRequests = len(statistics.rawResults)

	// is successful
	if workResult.StatusCode() > 199 && workResult.StatusCode() < 400 {
		statistics.numberOfSuccessfulRequests += 1
	} else {
		statistics.numberOfUnsuccessfulRequests += 1
	}

	// number of requests by status code
	statistics.numberOfRequestsByStatusCode[workResult.StatusCode()] += 1

	// number of requests by content type
	statistics.numberOfRequestsByContentType[workResult.ContentType()] += 1

	// update the total duration
	responseTime := workResult.EndTime().Sub(workResult.StartTime())
	statistics.totalResponseTime += responseTime

	// size
	statistics.totalSizeInBytes += workResult.Size()
	averageSizeInBytes := statistics.totalSizeInBytes / statistics.numberOfRequests

	// average response time
	averageResponseTime := time.Duration(statistics.totalResponseTime.Nanoseconds() / int64(statistics.numberOfRequests))

	// number of requests per second
	requestsPerSecond := statistics.endTime.Sub(statistics.startTime).Seconds() / float64(statistics.numberOfRequests)

	// create a snapshot
	snapShot := SnapShot{
		// times
		timestamp:           workResult.EndTime(),
		averageResponseTime: averageResponseTime,

		// counters
		totalNumberOfRequests:         statistics.numberOfRequests,
		numberOfSuccessfulRequests:    statistics.numberOfSuccessfulRequests,
		numberOfUnsuccessfulRequests:  statistics.numberOfUnsuccessfulRequests,
		numberOfRequestsPerSecond:     requestsPerSecond,
		numberOfRequestsByStatusCode:  statistics.numberOfRequestsByStatusCode,
		numberOfRequestsByContentType: statistics.numberOfRequestsByContentType,

		// size
		totalSizeInBytes:   statistics.totalSizeInBytes,
		averageSizeInBytes: averageSizeInBytes,
	}

	statistics.snapShots = append(statistics.snapShots, snapShot)

	return snapShot
}
