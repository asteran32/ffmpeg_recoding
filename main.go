package main

import (
	"gcs/record"
	"log"
	"sync"
)

var mutex = &sync.Mutex{}

func main() {
	// Use all CPU (4)
	// runtime.GOMAXPROCS(runtime.NumCPU() - 2)

	// Setting ffmpeg recoding
	f, err := record.InitRecording()
	if err != nil {
		log.Fatal(err)
	}
	if err = f.CreateDir(); err != nil {
		log.Fatal(err)
	}

	// // Setting google cloud storage
	// if err := cloud.InitStorage(); err != nil {
	// 	log.Fatal(err)
	// }
	// // Time ticker
	// recordTicker := time.NewTicker(time.Minute * 1)
	// uploadTicker := time.NewTicker(time.Minute * 5)

	done := make(chan bool, 1)
	go func() {
		for {
			f.Recording(f.Cam[0]) //for recording
			// cloud.UploadFiles("cam01/")
		}
	}()
	<-done
}
