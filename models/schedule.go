package models

type Schedule struct {
	ID          int    `json:"id"`
	LineID      int    `json:"line_id"`
	LineName    string `json:"line_name"`
	Description string `json:"description"`
	OnTime      string `json:"on_time"`
	OffTime     string `json:"off_time"`
	Scheduled   bool   `json:"scheduled"`
}
