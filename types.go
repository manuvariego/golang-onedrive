package main

type DriveItem struct {
	Name     string `json:"name"`
	IsFolder bool   `json:"folder,omitempty"`
	ID       string `json:"id"`
}
