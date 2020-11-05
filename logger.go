package main

import (
	"log"
)

func logResult(logger *log.Logger, workResult WorkResult) {
	logger.Println(workResult.String())
}
