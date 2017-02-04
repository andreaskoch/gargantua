package main

import (
	"fmt"
	"net/url"
)

type WorkRequest struct {
	URL     url.URL
	Execute func(workerID, numberOfWorkers int)
}

var WorkerQueue chan chan WorkRequest

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var workers []Worker

func StartDispatcher(numberOfWorkers int, stop chan bool) chan bool {

	allWorkersHaveStopped := make(chan bool)

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, numberOfWorkers)

	// Now, create all of our workers.
	for i := 0; i < numberOfWorkers; i++ {
		worker := NewWorker(i+1, WorkerQueue)
		workers = append(workers, worker)
		worker.Start()
	}

	go func() {
		for {
			select {

			case <-stop:
				fmt.Println("Receive the stop signal 2")
				for _, worker := range workers {
					worker.Stop()
				}

				allWorkersHaveStopped <- true
				return

			case work, ok := <-WorkQueue:
				if !ok {
					return
				}

				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()

	return allWorkersHaveStopped
}
