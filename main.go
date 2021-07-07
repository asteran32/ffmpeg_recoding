package main

import (
	"gcs/cloud"
	"gcs/record"
	"log"
	"sync"
	"time"
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
	if err := cloud.InitStorage(); err != nil {
		log.Fatal(err)
	}
	// // Time ticker
	recordTicker := time.NewTicker(time.Minute * 1)
	uploadTicker := time.NewTicker(time.Minute * 5)

	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done: // first
				break

			case t := <-recordTicker.C:
				timestr := t.Format("20060102-1504")
				for _, cam := range f.Cam {
					go f.RecordCam_single(timestr, cam)
				}
			case <-uploadTicker.C:
				for _, cam := range f.Cam {
					go cloud.UploadFiles(cam.Path)
				}

			default:
				continue
			}
		}
	}()
	<-done
}
