package controller

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func InitLogging() {
	// Logging
	y, month, d := time.Now().Date()
	logName := strconv.Itoa(y) + "_" + month.String() + "_" + strconv.Itoa(d)
	f, err := os.OpenFile("logs/log_"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	//Log to file:
	//log.SetOutput(f)

	//Log to console:
	log.SetOutput(os.Stdout)

	//Log both ways
	//log.SetOutput(io.MultiWriter(f, os.Stdout))

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("********APP STARTED********")
}
