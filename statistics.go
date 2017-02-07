package main

import (
	"fmt"
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

type Snapshot struct {

	// time
	timestamp           time.Time
	timeSinceStart      time.Duration
	averageResponseTime time.Duration

	// counters
	numberOfWorkers              int
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

func (snapshot Snapshot) Timestamp() time.Time {
	return snapshot.timestamp
}

func (snapshot Snapshot) NumberOfWorkers() int {
	return snapshot.numberOfWorkers
}

func (snapshot Snapshot) NumberOfErrors() int {
	return snapshot.numberOfUnsuccessfulRequests
}

func (snapshot Snapshot) TotalNumberOfRequests() int {
	return snapshot.totalNumberOfRequests
}

func (snapshot Snapshot) TotalSizeInBytes() int {
	return snapshot.totalSizeInBytes
}

func (snapshot Snapshot) AverageSizeInBytes() int {
	return snapshot.averageSizeInBytes
}

func (snapshot Snapshot) AverageResponseTime() time.Duration {
	return snapshot.averageResponseTime
}

func (snapshot Snapshot) RequestsPerSecond() float64 {
	return snapshot.numberOfRequestsPerSecond
}

type Statistics struct {
	lock sync.RWMutex

	rawResults  []WorkResult
	snapShots   []Snapshot
	logMessages []string

	startTime time.Time
	endTime   time.Time

	totalResponseTime time.Duration

	numberOfWorkers               int
	numberOfRequests              int
	numberOfSuccessfulRequests    int
	numberOfUnsuccessfulRequests  int
	numberOfRequestsByStatusCode  map[int]int
	numberOfRequestsByContentType map[string]int

	totalSizeInBytes int
}

func (statistics *Statistics) Add(workResult WorkResult) Snapshot {

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

	// number of workers
	statistics.numberOfWorkers = workResult.NumberOfWorkers()

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
	requestsPerSecond := float64(statistics.numberOfRequests) / statistics.endTime.Sub(statistics.startTime).Seconds()

	// log messages
	statistics.logMessages = append(statistics.logMessages, workResult.String())

	// create a snapshot
	snapShot := Snapshot{
		// times
		timestamp:           workResult.EndTime(),
		averageResponseTime: averageResponseTime,

		// counters
		numberOfWorkers:               statistics.numberOfWorkers,
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

func (statistics *Statistics) LastSnapshot() Snapshot {
	statistics.lock.RLock()
	defer statistics.lock.RUnlock()

	lastSnapshotIndex := len(statistics.snapShots) - 1
	if lastSnapshotIndex < 0 {
		return Snapshot{}
	}

	return statistics.snapShots[lastSnapshotIndex]
}

func (statistics *Statistics) LastLogMessages(count int) []string {

	statistics.lock.RLock()
	defer statistics.lock.RUnlock()

	messages, err := getLatestLogMessages(statistics.logMessages, count)
	if err != nil {
		panic(err)
	}

	return messages

}

func getLatestLogMessages(messages []string, count int) ([]string, error) {

	if count < 0 {
		return nil, fmt.Errorf("The count cannot be negative")
	}

	numberOfMessges := len(messages)

	if count == numberOfMessges {
		return messages, nil
	}

	if count < numberOfMessges {
		return messages[numberOfMessges-count:], nil
	}

	if count > numberOfMessges {
		fillLines := make([]string, count-numberOfMessges)
		return append(fillLines, messages...), nil
	}

	panic("Unreachable")
}
