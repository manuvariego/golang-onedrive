package main

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string `json:"name"`
	IsFolder    *bool  `json:"folder,omitifempty"`
	ID          string `json:"id"`
	DownloadUrl string `json:"@microsoft.graph.downloadUrl"`
}
