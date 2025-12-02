package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	models "github.com/pavelc4/ytdpl-api-go/internal/models"
)

type YTDLPService struct {
	cookiePath string
}

func NewYTDLPService(cookiePath string) *YTDLPService {
	return &YTDLPService{
		cookiePath: cookiePath,
	}
}

func (s *YTDLPService) GetDownloadURLs(url string) (*models.VideoURL, error) {
	args := []string{"-g", "--no-warnings"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to extract URLs: %w", err)
	}

	urls := strings.Split(strings.TrimSpace(string(output)), "\n")

	result := &models.VideoURL{
		VideoURL: urls[0],
	}

	if len(urls) > 1 {
		result.AudioURL = urls[1]
	}

	return result, nil
}

func (s *YTDLPService) GetVideoInfo(url string) (*models.VideoInfo, error) {
	args := []string{"-J", "--no-warnings"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to extract info: %w", err)
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

	return info, nil
}

func (s *YTDLPService) GetFormats(url string) (*models.FormatsResponse, error) {
	args := []string{"-J", "--no-warnings"}

	if s.cookiePath != "" {
		args = append(args, "--cookies", s.cookiePath)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to extract formats: %w", err)
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
