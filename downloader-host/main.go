package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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

func writeCookiesToFile(cookies []CookiePayload) (string, error) {
	tmpFile, err := os.CreateTemp("", "cookies-*.txt")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Simple Netscape-like format for yt-dlp
	fmt.Fprintf(tmpFile, "# Netscape HTTP Cookie File\n")
	for _, c := range cookies {
		fmt.Fprintf(tmpFile, "%s\tTRUE\t%s\tTRUE\t0\t%s\t%s\n", c.Domain, c.Path, c.Name, c.Value)
	}

	return tmpFile.Name(), nil
}

func worker(id int, jobs <-chan DownloadTask) {
	for task := range jobs {
		path := getDownloadPath(task.Site)
		log.Printf("Worker %d: starting download for %s\n", id, task.URL)

		cookiePath, err := writeCookiesToFile(task.Cookies)
		if err != nil {
			log.Printf("Worker %d: error writing cookies: %v\n", id, err)
			continue
		}

		cmd := exec.Command("yt-dlp", "--cookies", cookiePath, "-o", filepath.Join(path, "%(title)s.%(ext)s"), task.URL)
		output, err := cmd.CombinedOutput()
		os.Remove(cookiePath) // cleanup

		if err != nil {
			log.Printf("Worker %d: download failed: %v\nOutput: %s\n", id, err, string(output))
		} else {
			log.Printf("Worker %d: download complete: %s\n", id, task.URL)
		}
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
