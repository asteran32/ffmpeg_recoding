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

func (f *Ffmpeg) Recording(c Camera) {
	defer func() {
		mutex.Unlock()
	}()
	mutex.Lock()
	cmd := exec.Command("ffmpeg", "-rtsp_transport", "tcp", "-i", c.Url,
		"-c:v", "copy", "-threads", "0", "-an",
		"-f", "segment", "-strftime", "1", "-segment_time", "00:10:00",
		"-segment_atclocktime", "1", "-segment_clocktime_offset", "30", "-reset_timestamps", "1",
		"-segment_format", "mp4", c.Name+"-%Y-%m-%d_%H-%M.mp4")

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
