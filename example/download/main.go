package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	addr      = "127.0.0.1:1234"
	chunkSize = 64 * 1024         // 64 KB
	totalSize = 100 * 1024 * 1024 // 100 MB
)

func main() {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Download requested from %s", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=large.bin")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", totalSize))

		chunk := make([]byte, chunkSize)
		for i := range chunk {
			chunk[i] = byte(i)
		}

		remaining := totalSize
		for remaining > 0 {
			n := chunkSize
			if remaining < chunkSize {
				n = remaining
			}

			if _, err := w.Write(chunk[:n]); err != nil {
				log.Printf("client disconnected: %v", err)
				return
			}

			remaining -= n
		}

		log.Println("Download complete")
	})

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
