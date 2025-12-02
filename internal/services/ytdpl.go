package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"time"

	"github.com/patrickmn/go-cache"
	models "github.com/pavelc4/ytdpl-api-go/internal/models"
)

type YTDLPService struct {
	cookiePath string
	cache      *cache.Cache
	semaphore  chan struct{}
}

func NewYTDLPService(cookiePath string) *YTDLPService {
	return &YTDLPService{
		cookiePath: cookiePath,
		cache:      cache.New(15*time.Minute, 30*time.Minute),
		semaphore:  make(chan struct{}, 10), // Limit to 10 concurrent processes
	}
}

func (s *YTDLPService) GetDownloadURLs(url string) (*models.VideoURL, error) {
	// Check cache
	cacheKey := "dl_" + url
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*models.VideoURL), nil
	}

	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	args := []string{"-g", "--no-warnings", "--no-cache-dir", "--no-playlist"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to extract URLs: %w (output: %s)", err, string(output))
	}

	urls := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs found")
	}

	result := &models.VideoURL{
		VideoURL: urls[0],
	}

	if len(urls) > 1 {
		result.AudioURL = urls[1]
	}

	s.cache.Set(cacheKey, result, cache.DefaultExpiration)

	return result, nil
}

func (s *YTDLPService) GetVideoInfo(url string) (*models.VideoInfo, error) {
	// Check cache
	cacheKey := "info_" + url
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*models.VideoInfo), nil
	}

	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	args := []string{"-J", "--no-warnings", "--no-cache-dir"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to extract info: %w (output: %s)", err, string(output))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	info := &models.VideoInfo{
		ID:          getString(data, "id"),
		Title:       getString(data, "title"),
		Duration:    getInt(data, "duration"),
		Thumbnail:   getString(data, "thumbnail"),
		Description: getString(data, "description"),
		Uploader:    getString(data, "uploader"),
		ViewCount:   getInt(data, "view_count"),
		UploadDate:  getString(data, "upload_date"),
	}

	s.cache.Set(cacheKey, info, cache.DefaultExpiration)

	return info, nil
}

func (s *YTDLPService) GetFormats(url string) (*models.FormatsResponse, error) {
	// Check cache
	cacheKey := "fmt_" + url
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*models.FormatsResponse), nil
	}

	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	args := []string{"-J", "--no-warnings", "--no-cache-dir"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to extract formats: %w (output: %s)", err, string(output))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	response := &models.FormatsResponse{
		VideoID: getString(data, "id"),
		Formats: []models.VideoFormat{},
	}

	if fmts, ok := data["formats"].([]interface{}); ok {
		for _, f := range fmts {
			if fm, ok := f.(map[string]interface{}); ok {
				format := models.VideoFormat{
					FormatID:   getString(fm, "format_id"),
					Ext:        getString(fm, "ext"),
					Resolution: getString(fm, "resolution"),
					Quality:    getString(fm, "format_note"),
					Filesize:   getInt64(fm, "filesize"),
					FPS:        getInt(fm, "fps"),
					VCodec:     getString(fm, "vcodec"),
					ACodec:     getString(fm, "acodec"),
				}
				response.Formats = append(response.Formats, format)
			}
		}
	}

	s.cache.Set(cacheKey, response, cache.DefaultExpiration)

	return response, nil
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(float64); ok {
		return int64(val)
	}
	return 0
}
