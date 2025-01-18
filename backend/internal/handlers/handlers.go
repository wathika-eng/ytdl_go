package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"
	"videoDownloader/internal/lib"
	"videoDownloader/internal/manager"
)

var validURL = lib.IsValidURL

var downloadManager = &manager.DownloadManager{
	Downloads: make(map[string]*manager.Download),
}

func HandleSaveToDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" || !validURL(id) {
		http.Error(w, "Invalid download ID", http.StatusBadRequest)
		return
	}

	downloadManager.Mu.RLock()
	download, exists := downloadManager.Downloads[id]
	downloadManager.Mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	if download.Status.Status != "completed" {
		http.Error(w, "Download not completed", http.StatusBadRequest)
		return
	}

	fileName := lib.GetFileName(download.Status.URL)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "yt-dlp", "-o", "-", download.Status.URL)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to create stdout pipe: %v", err)
		http.Error(w, "Failed to stream video", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start yt-dlp: %v", err)
		http.Error(w, "Failed to stream video", http.StatusInternalServerError)
		return
	}
	defer cmd.Wait()

	_, err = io.Copy(w, stdout)
	if err != nil {
		log.Printf("Failed to stream video to client: %v", err)
		return
	}

	log.Printf("Download %s saved successfully.", id)
}

func HandleStartDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type, expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req manager.DownloadRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	log.Printf("Request body: %s", string(body))

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Starting download for URL: %s", req.URL)
	status, err := downloadManager.StartDownload(req.URL)
	if err != nil {
		log.Printf("Failed to start download: %v", err)
		http.Error(w, "Failed to start download: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Download started with ID: %s", status.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func HandlePauseDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Download ID required", http.StatusBadRequest)
		return
	}

	downloadManager.Mu.RLock()
	download, exists := downloadManager.Downloads[id]
	downloadManager.Mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	download.Mu.Lock()
	if download.Status.Status == "downloading" {
		close(download.Cancel)
		download.Status.Status = "paused"
	}
	download.Mu.Unlock()

	json.NewEncoder(w).Encode(download.Status)
}

func handleResumeDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Download ID required", http.StatusBadRequest)
		return
	}

	downloadManager.Mu.RLock()
	download, exists := downloadManager.Downloads[id]
	downloadManager.Mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	download.Mu.Lock()
	if download.Status.Status == "paused" {
		download.Cancel = make(chan struct{})
		go downloadManager.ProcessDownload(download)
	}
	download.Mu.Unlock()

	json.NewEncoder(w).Encode(download.Status)
}

func HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Download ID required", http.StatusBadRequest)
		return
	}

	downloadManager.Mu.RLock()
	download, exists := downloadManager.Downloads[id]
	downloadManager.Mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	download.Mu.RLock()
	status := manager.DownloadStatus{
		ID:        download.Status.ID,
		URL:       download.Status.URL,
		Progress:  download.Status.Progress,
		Status:    download.Status.Status,
		Error:     download.Status.Error,
		CreatedAt: download.Status.CreatedAt,
	}
	download.Mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func HandleListDownloads(w http.ResponseWriter, r *http.Request) {
	downloadManager.Mu.RLock()
	statuses := make([]*manager.DownloadStatus, 0, len(downloadManager.Downloads))
	for _, download := range downloadManager.Downloads {
		statuses = append(statuses, download.Status)
	}
	downloadManager.Mu.RUnlock()

	json.NewEncoder(w).Encode(statuses)
}
