package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type DownloadTask struct {
	URL     string          `json:"url"`
	Site    string          `json:"site"`
	Cookies []CookiePayload `json:"cookies"`
}

type CookiePayload struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func readMessage(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeMessage(w io.Writer, payload []byte) error {
	length := uint32(len(payload))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func getDownloadPath(site string) string {
	basePath := os.Getenv("BASE_DOWNLOAD_PATH")
	if basePath == "" {
		basePath = "downloads"
	}
	targetPath := filepath.Join(basePath, site)
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		log.Printf("Error creating directory %s: %v\n", targetPath, err)
	}
	return targetPath
}

func worker(id int, jobs <-chan DownloadTask) {
	for task := range jobs {
		path := getDownloadPath(task.Site)
		log.Printf("Worker %d: processing task for %s in %s\n", id, task.URL, path)
		// Task 3.5 will call yt-dlp here
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	numWorkersStr := os.Getenv("MAX_CONCURRENT_DOWNLOADS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil || numWorkers <= 0 {
		numWorkers = 3
	}

	jobs := make(chan DownloadTask, 100)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs)
	}

	for {
		payload, err := readMessage(os.Stdin)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Read error:", err)
		}

		var task DownloadTask
		if err := json.Unmarshal(payload, &task); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		jobs <- task

		resp := Response{Status: "queued", Message: "Download task added to queue"}
		respJSON, _ := json.Marshal(resp)
		if err := writeMessage(os.Stdout, respJSON); err != nil {
			log.Fatal("Write error:", err)
		}
	}
}
