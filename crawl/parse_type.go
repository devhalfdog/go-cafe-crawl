package crawl

/* Query struct */
type naverVideoData struct {
	Data struct {
		VideoType string `json:"videoType"`
		Vid       string `json:"vid"`
		Inkey     string `json:"inkey"`
	} `json:"data,omitempty"`
}
