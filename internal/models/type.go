package models

type YTDLPOutput struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Duration    float64       `json:"duration"`
	Thumbnail   string        `json:"thumbnail"`
	Description string        `json:"description"`
	Uploader    string        `json:"uploader"`
	ViewCount   int           `json:"view_count"`
	UploadDate  string        `json:"upload_date"`
	Formats     []VideoFormat `json:"formats"`
}

type VideoURL struct {
	VideoURL string `json:"video_url"`
	AudioURL string `json:"audio_url,omitempty"`
}

type VideoInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description,omitempty"`
	Uploader    string `json:"uploader"`
	ViewCount   int    `json:"view_count"`
	UploadDate  string `json:"upload_date,omitempty"`
}

type VideoFormat struct {
	FormatID   string  `json:"format_id"`
	Ext        string  `json:"ext"`
	Resolution string  `json:"resolution"`
	Quality    string  `json:"quality"`
	Filesize   int64   `json:"filesize"`
	FPS        float64 `json:"fps,omitempty"`
	VCodec     string  `json:"vcodec,omitempty"`
	ACodec     string  `json:"acodec,omitempty"`
}

type FormatsResponse struct {
	VideoID string        `json:"video_id"`
	Formats []VideoFormat `json:"formats"`
}
