package record

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

var mutex = &sync.Mutex{}

type Ffmpeg struct {
	FfmpegBin string   `json:"ffmpeg_bin"`
	Cam       []Camera `json:"camera"`
}

type Camera struct {
	Framerate string `json:"framerate"`
	Name      string `json:"name"`
	Codec     string `json:"codec"`
	Time      string `json:"time"`
	Path      string `json:"path"`
	Url       string `json:"url"`
}

func InitRecording() (Ffmpeg, error) {
	var ff Ffmpeg
	env, err := os.Open("cam.json")
	if err != nil {
		return ff, fmt.Errorf("Failed to open jsonfile: %v", err)
	}
	byteValue, _ := ioutil.ReadAll(env)
	if err = json.Unmarshal(byteValue, &ff); err != nil {
		return ff, fmt.Errorf("Faild to unmarshal: %v", err)
	}
	fmt.Println("Success to set Recording Env")
	return ff, nil
}

func (f *Ffmpeg) CreateDir() error {
	for _, c := range f.Cam { //key, value
		if _, err := os.Stat(c.Path); os.IsNotExist(err) {
			if err := os.MkdirAll(c.Path, 0755); err != nil {
				return fmt.Errorf("os.Mkdir: %v", err)
			}
		}
	}
	fmt.Println("Success to create video folders")
	return nil
}

func (f *Ffmpeg) RecordCam_multi(output []string) {
	defer func() {
		mutex.Unlock()
	}()
	mutex.Lock()

	fmt.Printf("Start to record videos : %v \n", output)

	// ffmpeg RTSP Command line only 3 camera
	cmd := exec.Command("ffmpeg", "-rtsp_transport", "tcp",
		"-i", f.Cam[0].Url, "-i", f.Cam[1].Url, "-i", f.Cam[2].Url,
		"-vcodec", "libx264", "-an", "-bufsize", "3M",
		"-map", "0:0", "-t", f.Cam[0].Time, output[0],
		"-map", "1:0", "-t", f.Cam[0].Time, output[1])

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	cmd.Start()
	oneByte := make([]byte, 1)
	for {
		_, err := stderr.Read(oneByte)
		if err != nil {
			break
		}
		fmt.Printf("%v", string(oneByte[0]))
	}

	cmd.Wait()

}

func (f *Ffmpeg) RecordCam_single(time string, c Camera) {
	defer func() {
		mutex.Unlock()
	}()

	mutex.Lock()
	output := c.Path + c.Name + "_" + time + ".mp4"

	cmd := exec.Command("ffmpeg", "-rtsp_transport", "tcp", "-i", c.Url,
		"-vcodec", c.Codec, "-an", "-max_delay", "1M", "-bufsize", "1M", "-t", c.Time, output)

	fmt.Printf("Start to record video: %v \n", output)

	cmd.Start()
}
