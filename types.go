package main

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string          `json:"name"`
	IsFolder    *bool           `json:"folder,omitifempty"`
	ID          string          `json:"id"`
	DownloadUrl string          `json:"@microsoft.graph.downloadUrl"`
	ParentData  ParentReference `json:"parentReference"`
}

type ParentReference struct {
	Path string `json:"path"`
	Name string `json:"name"`
}
