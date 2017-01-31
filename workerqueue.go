package main

import "net/url"

type WorkResult struct {
	Error   error
	Message string
}

type WorkRequest struct {
	URL     url.URL
	Execute func() WorkResult
}

var WorkerQueue chan chan WorkRequest

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

func StartDispatcher(nworkers int) chan WorkResult {

	results := make(chan WorkResult, 10)

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
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
