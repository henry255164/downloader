package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type DownloadTask struct {
	URL     string          `json:"url"`
	Site    string          `json:"site"`
	Cookies []CookiePayload `json:"cookies"`
}

type CookiePayload struct {
	Name           string  `json:"name"`
	Value          string  `json:"value"`
	Domain         string  `json:"domain"`
	Path           string  `json:"path"`
	Secure         bool    `json:"secure"`
	ExpirationDate float64 `json:"expirationDate"`
	HostOnly       bool    `json:"hostOnly"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var wg sync.WaitGroup

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
	if !filepath.IsAbs(basePath) {
		exePath, _ := os.Executable()
		basePath = filepath.Join(filepath.Dir(exePath), basePath)
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

	// Official Netscape format: 
	// domain  include_subdomains  path  secure  expiry  name  value
	fmt.Fprintf(tmpFile, "# Netscape HTTP Cookie File\n")
	for _, c := range cookies {
		// Rule: If domain starts with '.', flag is TRUE. If host-only, flag is FALSE.
		includeSubdomains := "TRUE"
		if c.HostOnly {
			includeSubdomains = "FALSE"
		} else if !strings.HasPrefix(c.Domain, ".") {
			// Ensure domain cookies start with a dot if not host-only
			// but we follow the flag provided by Chrome
		}

		secure := "FALSE"
		if c.Secure {
			secure = "TRUE"
		}

		// Use 0 if expiration date is missing (session cookie)
		expiry := int64(c.ExpirationDate)
		
		// Ensure name and value don't contain tabs that break the format
		name := strings.ReplaceAll(c.Name, "\t", " ")
		value := strings.ReplaceAll(c.Value, "\t", " ")

		fmt.Fprintf(tmpFile, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n", 
			c.Domain, includeSubdomains, c.Path, secure, expiry, name, value)
	}

	return tmpFile.Name(), nil
}

func worker(id int, jobs <-chan DownloadTask) {
	for task := range jobs {
		wg.Add(1)
		executeDownload(id, task)
		wg.Done()
	}
}

func executeDownload(id int, task DownloadTask) {
	path := getDownloadPath(task.Site)
	
	ytDlpPath := os.Getenv("YT_DLP_PATH")
	if ytDlpPath == "" {
		ytDlpPath = "yt-dlp"
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")

	cookiePath, err := writeCookiesToFile(task.Cookies)
	if err != nil {
		log.Printf("Worker %d: error writing cookies: %v\n", id, err)
		return
	}
	defer os.Remove(cookiePath)

	outputPath := filepath.Join(path, "%(title)s.%(ext)s")
	
	args := []string{"--newline", "--cookies", cookiePath, "-o", outputPath}
	if ffmpegPath != "" {
		args = append(args, "--ffmpeg-location", ffmpegPath)
	}
	args = append(args, task.URL)

	log.Printf("Worker %d: Executing [%s %v]\n", id, ytDlpPath, args)

	cmd := exec.Command(ytDlpPath, args...)
	
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("Worker %d: FAILED TO START: %v\n", id, err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Printf("W%d-OUT: %s\n", id, scanner.Text())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("W%d-ERR: %s\n", id, scanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Printf("Worker %d: Download finished with error: %v\n", id, err)
		exec.Command("osascript", "-e", fmt.Sprintf("display notification \"Download failed for %s\" with title \"Media Downloader\"", task.URL)).Run()
	} else {
		log.Printf("Worker %d: Download complete: %s\n", id, task.URL)
		exec.Command("osascript", "-e", fmt.Sprintf("display notification \"Download complete!\" subtitle \"%s\" with title \"Media Downloader\"", task.URL)).Run()
	}
}

func main() {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	logFile := filepath.Join(exeDir, "app.log")
	
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(f)
		defer f.Close()
	}

	godotenv.Load(filepath.Join(exeDir, ".env"))
	log.Printf("Service started. PID: %d, YT_DLP: %s, FFMPEG: %s\n", os.Getpid(), os.Getenv("YT_DLP_PATH"), os.Getenv("FFMPEG_PATH"))

	numWorkersStr := os.Getenv("MAX_CONCURRENT_DOWNLOADS")
	numWorkers, _ := strconv.Atoi(numWorkersStr)
	if numWorkers <= 0 {
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
				log.Println("Stdin closed (EOF), waiting for workers to finish...")
				break
			}
			log.Printf("Read error: %v\n", err)
			break
		}

		var task DownloadTask
		if err := json.Unmarshal(payload, &task); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		jobs <- task
		writeMessage(os.Stdout, []byte(`{"status":"queued"}`))
	}

	wg.Wait()
	log.Println("Service exiting normally.")
}
