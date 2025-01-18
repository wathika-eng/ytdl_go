package main

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"
	"videoDownloader/internal/handlers"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	templatesSubFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			indexFile, err := templatesSubFS.Open("index.html")
			if err != nil {
				http.Error(w, "Index file not found", http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			content, err := io.ReadAll(indexFile)
			if err != nil {
				http.Error(w, "Failed to read index file", http.StatusInternalServerError)
				return
			}
			http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(content))
			return
		}

		http.FileServer(http.FS(templatesSubFS)).ServeHTTP(w, r)
	})

	http.HandleFunc("/api/download/test", handlers.TestAPI)
	http.HandleFunc("/api/download/start", handlers.HandleStartDownload)
	http.HandleFunc("/api/download/pause", handlers.HandlePauseDownload)
	http.HandleFunc("/api/download/resume", handlers.HandlePauseDownload)
	http.HandleFunc("/api/download/status", handlers.HandleGetStatus)
	http.HandleFunc("/api/downloads", handlers.HandleListDownloads)
	http.HandleFunc("/api/download/save", handlers.HandleSaveToDevice)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
