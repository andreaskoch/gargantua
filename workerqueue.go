package main

import "fmt"

type WorkResult struct {
	Message string
}

type WorkRequest struct {
	Name    string
	Execute func() WorkResult
}

var WorkerQueue chan chan WorkRequest

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

func StartDispatcher(nworkers int, numberOfWorkRequets int) chan WorkResult {

	results := make(chan WorkResult, numberOfWorkRequets)

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
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
