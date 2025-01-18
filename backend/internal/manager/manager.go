package manager

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
	"videoDownloader/internal/lib"
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
	Downloads map[string]*Download
	Mu        sync.RWMutex
}

type Download struct {
	Status *DownloadStatus
	Cmd    *exec.Cmd
	Cancel chan struct{}
	Mu     sync.RWMutex
}

var validURL = lib.IsValidURL

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (dm *DownloadManager) StartDownload(url string) (*DownloadStatus, error) {
	if !validURL(url) {
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
		Cancel: make(chan struct{}),
	}

	dm.Mu.Lock()
	dm.Downloads[id] = download
	dm.Mu.Unlock()

	go dm.ProcessDownload(download)

	return status, nil
}

func ParseProgress(line string) float64 {
	matches := lib.ProgressRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		if progress, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return progress
		}
	}
	return -1
}

func (dm *DownloadManager) ProcessDownload(download *Download) {
	download.Mu.Lock()
	download.Status.Status = "downloading"
	download.Status.Progress = 0
	download.Mu.Unlock()

	cmd := exec.Command("yt-dlp",
		"-o", "-",
		"--newline",
		"--progress",
		"--progress-template", "%(progress._percent_str)s",
		download.Status.URL)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		download.Mu.Lock()
		download.Status.Status = "error"
		download.Status.Error = err.Error()
		download.Mu.Unlock()
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		download.Mu.Lock()
		download.Status.Status = "error"
		download.Status.Error = err.Error()
		download.Mu.Unlock()
		return
	}

	download.Cmd = cmd
	if err := cmd.Start(); err != nil {
		download.Mu.Lock()
		download.Status.Status = "error"
		download.Status.Error = err.Error()
		download.Mu.Unlock()
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case <-download.Cancel:
				cmd.Process.Kill()
				return
			default:
				text := scanner.Text()
				progress := ParseProgress(text)
				if progress >= 0 {
					download.Mu.Lock()
					download.Status.Progress = progress
					download.Mu.Unlock()
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
		download.Mu.Lock()
		if download.Status.Status != "paused" {
			download.Status.Status = "error"
			download.Status.Error = err.Error()
		}
		download.Mu.Unlock()
		return
	}

	download.Mu.Lock()
	if download.Status.Status != "paused" {
		download.Status.Status = "completed"
		download.Status.Progress = 100
	}
	download.Mu.Unlock()
}
