package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
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

// readMessage reads a native message from stdin (4-byte length prefix + payload)
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

// writeMessage writes a native message to stdout (4-byte length prefix + payload)
func writeMessage(w io.Writer, payload []byte) error {
	length := uint32(len(payload))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func main() {
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

		log.Printf("Received task: %s from site %s\n", task.URL, task.Site)

		// Respond back
		resp := Response{Status: "received", Message: "Download task queued"}
		respJSON, _ := json.Marshal(resp)
		if err := writeMessage(os.Stdout, respJSON); err != nil {
			log.Fatal("Write error:", err)
		}
	}
}
