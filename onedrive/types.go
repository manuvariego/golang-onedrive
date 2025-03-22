package onedrive

import (
	"net/http"
	"sync"
)

var oneDriveRoot = Directory{Name: "root"}
var mu sync.Mutex

type OneDriver interface {
	Pwd() Path
	Ls() ([]string, error)
	Cd(folder string) (Path, error)
}

type OneDriveClient struct {
	Client *http.Client
	Path   Path
}

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string `json:"name"`
	IsFolder    *bool  `json:"folder,omitempty"`
	ID          string `json:"id"`
	DownloadUrl string `json:"@microsoft.graph.downloadUrl"`
}

type Directory struct {
	Name     string      `json:"name"`
	Files    []File      `json:"files,omitempty"`
	Children []Directory `json:"folders,omitempty"`
}

type File struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DownloadUrl string `json:"download_url"`
}

type Path struct {
	CurrentPath string
}
