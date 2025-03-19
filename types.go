package main

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string          `json:"name"`
	IsFolder    *bool           `json:"folder,omitempty"`
	ID          string          `json:"id"`
	DownloadUrl string          `json:"@microsoft.graph.downloadUrl"`
	ParentData  ParentReference `json:"parentReference"`
}

type ParentReference struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type Path struct {
	CurrentPath string
	ParentPath  string
}
