package main

import (
	"fmt"
	"time"

	"github.com/gizak/termui"
)

func dashboard(startTime time.Time, stopTheCrawler chan bool) {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	var snapshots []Snapshot

	logWindow := termui.NewList()
	logWindow.ItemFgColor = termui.ColorYellow
	logWindow.BorderLabel = "Log"
	logWindow.Height = 22
	logWindow.Y = 0

	totalBytesDownloaded := termui.NewPar("")
	totalBytesDownloaded.Height = 3
	totalBytesDownloaded.TextFgColor = termui.ColorWhite
	totalBytesDownloaded.BorderLabel = "Data downloaded"
	totalBytesDownloaded.BorderFg = termui.ColorCyan

	totalNumberOfRequests := termui.NewPar("")
	totalNumberOfRequests.Height = 3
	totalNumberOfRequests.TextFgColor = termui.ColorWhite
	totalNumberOfRequests.BorderLabel = "URLs crawled"
	totalNumberOfRequests.BorderFg = termui.ColorCyan

	requestsPerSecond := termui.NewPar("")
	requestsPerSecond.Height = 3
	requestsPerSecond.TextFgColor = termui.ColorWhite
	requestsPerSecond.BorderLabel = "URLs/second"
	requestsPerSecond.BorderFg = termui.ColorCyan

	averageResponseTime := termui.NewPar("")
	averageResponseTime.Height = 3
	averageResponseTime.TextFgColor = termui.ColorWhite
	averageResponseTime.BorderLabel = "Average response time"
	averageResponseTime.BorderFg = termui.ColorCyan

	numberOfWorkers := termui.NewPar("")
	numberOfWorkers.Height = 3
	numberOfWorkers.TextFgColor = termui.ColorWhite
	numberOfWorkers.BorderLabel = "Number of workers"
	numberOfWorkers.BorderFg = termui.ColorCyan

	averageSizeInBytes := termui.NewPar("")
	averageSizeInBytes.Height = 3
	averageSizeInBytes.TextFgColor = termui.ColorWhite
	averageSizeInBytes.BorderLabel = "Average response size"
	averageSizeInBytes.BorderFg = termui.ColorCyan

	numberOfErrors := termui.NewPar("")
	numberOfErrors.Height = 3
	numberOfErrors.TextFgColor = termui.ColorWhite
	numberOfErrors.BorderLabel = "Number of 4xx errors"
	numberOfErrors.BorderFg = termui.ColorCyan

	elapsedTime := termui.NewPar("")
	elapsedTime.Height = 3
	elapsedTime.TextFgColor = termui.ColorWhite
	elapsedTime.BorderLabel = "Elapsed time"
	elapsedTime.BorderFg = termui.ColorCyan

	draw := func() {

		snapshot := stats.LastSnapshot()

		// capture the latest snapshot
		snapshots = append(snapshots, snapshot)

		// log messages
		logWindow.Items = stats.LastLogMessages(20)

		// total number of requests
		totalNumberOfRequests.Text = fmt.Sprintf("%d", snapshot.TotalNumberOfRequests())

		// total number of bytes downloaded
		totalBytesDownloaded.Text = formatBytes(snapshot.TotalSizeInBytes())

		// requests per second
		requestsPerSecond.Text = fmt.Sprintf("%.1f", snapshot.RequestsPerSecond())

		// average response time
		averageResponseTime.Text = fmt.Sprintf("%s", snapshot.AverageResponseTime())

		// number of workers
		numberOfWorkers.Text = fmt.Sprintf("%d", snapshot.NumberOfWorkers())

		// average request size
		averageSizeInBytes.Text = formatBytes(snapshot.AverageSizeInBytes())

		// number of errors
		numberOfErrors.Text = fmt.Sprintf("%d", snapshot.NumberOfErrors())

		// time since first snapshot
		timeSinceStart := time.Now().Sub(startTime)
		elapsedTime.Text = fmt.Sprintf("%s", timeSinceStart)

		termui.Render(termui.Body)
	}

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, logWindow),
		),
		termui.NewRow(
			termui.NewCol(3, 0, totalBytesDownloaded),
			termui.NewCol(3, 0, totalNumberOfRequests),
			termui.NewCol(3, 0, requestsPerSecond),
			termui.NewCol(3, 0, averageResponseTime),
		),
		termui.NewRow(
			termui.NewCol(3, 0, numberOfWorkers),
			termui.NewCol(3, 0, numberOfErrors),
			termui.NewCol(3, 0, averageSizeInBytes),
			termui.NewCol(3, 0, elapsedTime),
		),
	)

	termui.Body.Align()

	termui.Render(termui.Body)

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		stopTheCrawler <- true

		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		draw()
	})

	termui.Loop()
}
