package main

import "net/http"

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string `json:"name"`
	IsFolder    *bool  `json:"folder,omitempty"`
	ID          string `json:"id"`
	DownloadUrl string `json:"@microsoft.graph.downloadUrl"`
}

// type ParentReference struct {
// 	Path string `json:"path"`
// 	Name string `json:"name"`
// }

type Path struct {
	CurrentPath string
}

type OneDriveClient struct {
	Client *http.Client
	Path   Path
}

type OneDriver interface {
	Pwd() Path
	Ls() ([]string, error)
	Cd(folder string) (Path, error)
}
