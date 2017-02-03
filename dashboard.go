package main

import ui "github.com/gizak/termui"
import "math"

func dashboard() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	strs := []string{"[0] gizak/termui", "[1] editbox.go", "[2] iterrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}
	logWindow := ui.NewList()
	logWindow.Items = strs
	logWindow.ItemFgColor = ui.ColorYellow
	logWindow.BorderLabel = "Log"
	logWindow.Height = 14
	logWindow.Width = 54
	logWindow.Y = 0

	totalBytesDownloaded := ui.NewPar("0")
	totalBytesDownloaded.Width = 25
	totalBytesDownloaded.Height = 3
	totalBytesDownloaded.Y = 11
	totalBytesDownloaded.X = 56
	totalBytesDownloaded.TextFgColor = ui.ColorWhite
	totalBytesDownloaded.BorderLabel = "Bytes downloaded"
	totalBytesDownloaded.BorderFg = ui.ColorCyan

	totalNumberOfRequests := ui.NewPar("0")
	totalNumberOfRequests.Width = 25
	totalNumberOfRequests.Height = 3
	totalNumberOfRequests.Y = 14
	totalNumberOfRequests.X = 56
	totalNumberOfRequests.TextFgColor = ui.ColorWhite
	totalNumberOfRequests.BorderLabel = "Number of requests"
	totalNumberOfRequests.BorderFg = ui.ColorCyan

	averageSizeInBytes := ui.Sparkline{}
	averageSizeInBytes.Height = 1
	averageSizeInBytes.Title = "⟨size⟩"
	spdata := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
	averageSizeInBytes.Data = spdata
	averageSizeInBytes.LineColor = ui.ColorCyan
	averageSizeInBytes.TitleColor = ui.ColorWhite

	averageDuration := ui.Sparkline{}
	averageDuration.Height = 1
	averageDuration.Title = "⟨duration⟩"
	averageDuration.Data = spdata
	averageDuration.TitleColor = ui.ColorWhite
	averageDuration.LineColor = ui.ColorRed

	averages := ui.NewSparklines(averageSizeInBytes, averageDuration)
	averages.BorderLabel = "Averages"
	averages.Width = 25
	averages.Height = 7
	averages.Y = 18
	averages.X = 56

	sinps := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	requestsByStatusCode := ui.NewBarChart()
	bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	requestsByStatusCode.BorderLabel = "Requests by status code"
	requestsByStatusCode.Width = 26
	requestsByStatusCode.Height = 11
	requestsByStatusCode.X = 56
	requestsByStatusCode.Y = 0
	requestsByStatusCode.DataLabels = bclabels
	requestsByStatusCode.BarColor = ui.ColorGreen
	requestsByStatusCode.NumColor = ui.ColorBlack

	requestsPerSecondLineChart := ui.NewLineChart()
	requestsPerSecondLineChart.BorderLabel = "Requests per second"
	requestsPerSecondLineChart.Data = sinps
	requestsPerSecondLineChart.Width = 54
	requestsPerSecondLineChart.Height = 11
	requestsPerSecondLineChart.X = 0
	requestsPerSecondLineChart.Y = 14
	requestsPerSecondLineChart.AxesColor = ui.ColorWhite
	requestsPerSecondLineChart.LineColor = ui.ColorYellow | ui.AttrBold

	draw := func(t int) {

		logWindow.Items = strs[t%9:]
		averages.Lines[0].Data = spdata[:30+t%50]
		averages.Lines[1].Data = spdata[:35+t%50]
		requestsPerSecondLineChart.Data = sinps[2*t%220:]
		requestsByStatusCode.Data = bcdata[t/2%10:]
		ui.Render(logWindow, requestsByStatusCode, totalBytesDownloaded, totalNumberOfRequests, averages, requestsPerSecondLineChart)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})
	ui.Loop()
}
