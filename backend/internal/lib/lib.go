package lib

import (
	"net/url"
	"path"
	"regexp"
	"strings"
)

var ProgressRegex = regexp.MustCompile(`(\d+\.\d+)%`)

func IsValidURL(url string) bool {
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

func IsValidID(id string) bool {
	// Validate ID format (alphanumeric, 8-32 characters)
	match, _ := regexp.MatchString(`^[a-zA-Z0-9]{8,32}$`, id)
	return match
}

func GetFileName(rawURL string) string {
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
