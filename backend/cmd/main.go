package main

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed templates/*
var templatesFS embed.FS

var (
	progressRegex = regexp.MustCompile(`(\d+\.\d+)%`)
)

type DownloadRequest struct {
	URL string `json:"url"`
}

type DownloadStatus struct {
	ID        string  `json:"id"`
	URL       string  `json:"url"`
	Progress  float64 `json:"progress"`
	Status    string  `json:"status"`
	Error     string  `json:"error,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

type DownloadManager struct {
	downloads map[string]*Download
	mu        sync.RWMutex
}

type Download struct {
	Status *DownloadStatus
	cmd    *exec.Cmd
	cancel chan struct{}
	mu     sync.RWMutex
}

var downloadManager = &DownloadManager{
	downloads: make(map[string]*Download),
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func isValidURL(url string) bool {
	// Supported platforms: YouTube, TikTok, Instagram Reels, Twitter/X, Pornhub
	platformRegex := []*regexp.Regexp{
		regexp.MustCompile(`^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+`),                                    // YouTube
		regexp.MustCompile(`^(https?:\/\/)?(www\.)?(tiktok\.com\/@.+\/video\/\d+|vm\.tiktok\.com\/[A-Za-z0-9]+\/?)`), // TikTok
		regexp.MustCompile(`^(https?:\/\/)?(www\.)?instagram\.com\/reels?\/.+`),                                      // Instagram Reels
		regexp.MustCompile(`^(https?:\/\/)?(www\.)?(twitter\.com|x\.com)\/.+\/status\/\d+`),                          // Twitter/X
		regexp.MustCompile(`^(https?:\/\/)?(www\.)?pornhub\.com\/view_video\.php\?viewkey=[a-zA-Z0-9]+`),             // Pornhub
	}

	for _, regex := range platformRegex {
		if regex.MatchString(url) {
			return true
		}
	}
	return false
}

func (dm *DownloadManager) StartDownload(url string) (*DownloadStatus, error) {
	if !isValidURL(url) {
		return nil, fmt.Errorf("invalid URL: only YouTube, TikTok, Instagram Reels, Twitter/X, and Pornhub are supported")
	}

	id := generateID()
	status := &DownloadStatus{
		ID:        id,
		URL:       url,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	download := &Download{
		Status: status,
		cancel: make(chan struct{}),
	}

	dm.mu.Lock()
	dm.downloads[id] = download
	dm.mu.Unlock()

	go dm.processDownload(download)

	return status, nil
}

func parseProgress(line string) float64 {
	matches := progressRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		if progress, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return progress
		}
	}
	return -1
}

func (dm *DownloadManager) processDownload(download *Download) {
	download.mu.Lock()
	download.Status.Status = "downloading"
	download.Status.Progress = 0
	download.mu.Unlock()

	cmd := exec.Command("yt-dlp",
		"-o", "-",
		"--newline",
		"--progress",
		"--progress-template", "%(progress._percent_str)s",
		download.Status.URL)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		download.mu.Lock()
		download.Status.Status = "error"
		download.Status.Error = err.Error()
		download.mu.Unlock()
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		download.mu.Lock()
		download.Status.Status = "error"
		download.Status.Error = err.Error()
		download.mu.Unlock()
		return
	}

	download.cmd = cmd
	if err := cmd.Start(); err != nil {
		download.mu.Lock()
		download.Status.Status = "error"
		download.Status.Error = err.Error()
		download.mu.Unlock()
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case <-download.cancel:
				cmd.Process.Kill()
				return
			default:
				text := scanner.Text()
				progress := parseProgress(text)
				if progress >= 0 {
					download.mu.Lock()
					download.Status.Progress = progress
					download.mu.Unlock()
				}
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("Download %s stderr: %s", download.Status.ID, scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		download.mu.Lock()
		if download.Status.Status != "paused" {
			download.Status.Status = "error"
			download.Status.Error = err.Error()
		}
		download.mu.Unlock()
		return
	}

	download.mu.Lock()
	if download.Status.Status != "paused" {
		download.Status.Status = "completed"
		download.Status.Progress = 100
	}
	download.mu.Unlock()
}

func isValidID(id string) bool {
	// Validate ID format (alphanumeric, 8-32 characters)
	match, _ := regexp.MatchString(`^[a-zA-Z0-9]{8,32}$`, id)
	return match
}

func getFileName(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "downloaded_video.mp4"
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "" || fileName == "/" || fileName == "." {
		return "downloaded_video.mp4"
	}

	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "#")[0]

	if !strings.Contains(fileName, ".") {
		fileName += ".mp4"
	}

	return fileName
}
func handleSaveToDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" || !isValidID(id) {
		http.Error(w, "Invalid download ID", http.StatusBadRequest)
		return
	}

	downloadManager.mu.RLock()
	download, exists := downloadManager.downloads[id]
	downloadManager.mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	if download.Status.Status != "completed" {
		http.Error(w, "Download not completed", http.StatusBadRequest)
		return
	}

	fileName := getFileName(download.Status.URL)
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

func handleStartDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type, expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req DownloadRequest
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

func handlePauseDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Download ID required", http.StatusBadRequest)
		return
	}

	downloadManager.mu.RLock()
	download, exists := downloadManager.downloads[id]
	downloadManager.mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	download.mu.Lock()
	if download.Status.Status == "downloading" {
		close(download.cancel)
		download.Status.Status = "paused"
	}
	download.mu.Unlock()

	json.NewEncoder(w).Encode(download.Status)
}

func handleResumeDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Download ID required", http.StatusBadRequest)
		return
	}

	downloadManager.mu.RLock()
	download, exists := downloadManager.downloads[id]
	downloadManager.mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	download.mu.Lock()
	if download.Status.Status == "paused" {
		download.cancel = make(chan struct{})
		go downloadManager.processDownload(download)
	}
	download.mu.Unlock()

	json.NewEncoder(w).Encode(download.Status)
}

func handleGetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Download ID required", http.StatusBadRequest)
		return
	}

	downloadManager.mu.RLock()
	download, exists := downloadManager.downloads[id]
	downloadManager.mu.RUnlock()

	if !exists {
		http.Error(w, "Download not found", http.StatusNotFound)
		return
	}

	download.mu.RLock()
	status := DownloadStatus{
		ID:        download.Status.ID,
		URL:       download.Status.URL,
		Progress:  download.Status.Progress,
		Status:    download.Status.Status,
		Error:     download.Status.Error,
		CreatedAt: download.Status.CreatedAt,
	}
	download.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleListDownloads(w http.ResponseWriter, r *http.Request) {
	downloadManager.mu.RLock()
	statuses := make([]*DownloadStatus, 0, len(downloadManager.downloads))
	for _, download := range downloadManager.downloads {
		statuses = append(statuses, download.Status)
	}
	downloadManager.mu.RUnlock()

	json.NewEncoder(w).Encode(statuses)
}

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

	http.HandleFunc("/api/download/start", handleStartDownload)
	http.HandleFunc("/api/download/pause", handlePauseDownload)
	http.HandleFunc("/api/download/resume", handleResumeDownload)
	http.HandleFunc("/api/download/status", handleGetStatus)
	http.HandleFunc("/api/downloads", handleListDownloads)
	http.HandleFunc("/api/download/save", handleSaveToDevice)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
