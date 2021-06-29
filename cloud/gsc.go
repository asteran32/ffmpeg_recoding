package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"

	"cloud.google.com/go/storage"
)

var mutex = &sync.Mutex{}

var Gcs = GCS{}

type GCS struct {
	Projcet     string `json:"projcet"`
	Bucket      string `json:"bucket"`
	Credentials string `json:"credentials"`
}

// $env:GOOGLE_APPLICATION_CREDENTIALS='C:\Users\AAA\Desktop\Fine\record_ipcam\fine-dt-project-2021-06-11-3a741164c4ab.json'
// export GOOGLE_APPLICATION_CREDENTIALS="/home/godheeran/Fine_recoring_sys/fine-dt-project-2021-06-11-3a741164c4ab.json"
func InitStorage() error {
	env, err := os.Open("gcs.json")
	if err != nil {
		return fmt.Errorf("Failed to open jsonfile: %v", err)
	}
	byteValue, _ := ioutil.ReadAll(env)
	if err = json.Unmarshal(byteValue, &Gcs); err != nil {
		return fmt.Errorf("Faild to unmarshal: %v", err)
	}

	fmt.Println("Success to set google cloud storage environments")

	return nil
}

// Set google application credentials key path
func (g *GCS) SetCredentials() error {
	cmd := exec.Command("cmd", g.Credentials)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to set google cloud storage key")
		return err
	}
	fmt.Println("Successfully set google cloud storage key")
	return nil
}

// Upload video file on google cloud storage
func (g *GCS) UploadFile(filePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	if err = uploadToCloud(client, filePath, g.Bucket); err != nil {
		return err
	}

	fmt.Printf("Success to upload bolb : %v \n", filePath)

	return nil
}

// Upload video file on google cloud storage
func UploadFiles(dirPath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	dir, _ := ioutil.ReadDir(dirPath)

	for _, fpath := range dir {
		if err = uploadToCloud(client, dirPath+fpath.Name(), Gcs.Bucket); err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func uploadToCloud(client *storage.Client, filePath string, bucket string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Open local file
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(filePath).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	// Pubilc Setting
	acl := client.Bucket(bucket).Object(filePath).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}

	// Delete local files
	if err = deleteLocal(filePath); err != nil {
		return fmt.Errorf("os.Remove: %v", err)
	}

	fmt.Printf("Success to upload bolb : %v \n", filePath)
	return nil
}

func deleteLocal(output string) error {
	if _, err := os.Stat(output); os.IsNotExist(err) {
		return fmt.Errorf("%v is not exsited", output)
	}
	os.Remove(output)
	return nil
}
