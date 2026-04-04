package models

type Line struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	IsOn bool   `json:"is_on"`
}
