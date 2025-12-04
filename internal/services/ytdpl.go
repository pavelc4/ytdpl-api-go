package services

import (
	"context"
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

func (s *YTDLPService) GetDownloadURLs(ctx context.Context, url string) (*models.VideoURL, error) {
	cacheKey := "dl_" + url
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*models.VideoURL), nil
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	args := []string{"-g", "--no-warnings", "--no-cache-dir", "--no-playlist"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
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

func (s *YTDLPService) GetVideoInfo(ctx context.Context, url string) (*models.VideoInfo, error) {
	cacheKey := "info_" + url
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*models.VideoInfo), nil
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	args := []string{"-J", "--no-warnings", "--no-cache-dir"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to extract info: %w (output: %s)", err, string(output))
	}

	var data models.YTDLPOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	info := &models.VideoInfo{
		ID:          data.ID,
		Title:       data.Title,
		Duration:    int(data.Duration),
		Thumbnail:   data.Thumbnail,
		Description: data.Description,
		Uploader:    data.Uploader,
		ViewCount:   data.ViewCount,
		UploadDate:  data.UploadDate,
	}

	s.cache.Set(cacheKey, info, cache.DefaultExpiration)

	return info, nil
}

func (s *YTDLPService) GetFormats(ctx context.Context, url string) (*models.FormatsResponse, error) {
	cacheKey := "fmt_" + url
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*models.FormatsResponse), nil
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	args := []string{
		"-J",
		"--no-playlist",
		"--no-warnings",
		"--no-cache-dir",
	}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to extract formats: %w (output: %s)", err, string(output))
	}

	var data models.YTDLPOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	response := &models.FormatsResponse{
		VideoID: data.ID,
		Formats: data.Formats,
	}

	s.cache.Set(cacheKey, response, cache.DefaultExpiration)

	return response, nil
}

func (s *YTDLPService) DownloadToFile(ctx context.Context, url, outputPath, quality string) error {
	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return ctx.Err()
	}

	format := "bestvideo+bestaudio/best"
	if quality == "720p" {
		format = "bestvideo[height<=720]+bestaudio/best[height<=720]"
	} else if quality == "1080p" {
		format = "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
	}

	args := []string{
		"-f", format,
		"--merge-output-format", "mp4",
		"--no-playlist",
		"--no-warnings",
		"--no-cache-dir",
		"-o", outputPath,
	}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download video: %w (output: %s)", err, string(output))
	}

	return nil
}
