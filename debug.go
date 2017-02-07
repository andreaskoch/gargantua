package main

import "fmt"

var errorMessages []string

var debugf func(format string, a ...interface{})

func init() {
	debugf = func(format string, a ...interface{}) {
		// default
	}
}

func consoleDebug(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func dashboardDebug(format string, a ...interface{}) {
	latestMesasges, err := getLatestLogMessages(errorMessages, 4)
	if err != nil {
		panic(err)
	}

	errorMessages = append(latestMesasges, fmt.Sprintf(format, a...))
}
