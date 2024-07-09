package models

type Panel struct {
	ID int `json:"id"`
}

type Dashboard struct {
	Title string `json:"title"`
	Panels []Panel `json:"panels"`
}