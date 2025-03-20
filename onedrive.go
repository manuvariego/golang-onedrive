package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func NewRootItem(sharepoint string) *Item {
	return &Item{
		path:   fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drives/%s", sharepoint),
		parent: nil,
	}
}

func directoryExists(items []Item, cmd string) bool {
	for _, item := range items {
		if item.Name == cmd && item.IsFolder != nil {
			return true
		}
	}
	return false
}

func (i *Item) ChangeDirectory(dest string) *Item {
	if dest == ".." {
		if i.parent != nil {
			return i.parent
		}
		return i
	}
	return &Item{
		path:   i.path + "/" + dest,
		parent: i,
	}
}

func (i *Item) Pwd() string {
	return i.path
}

// Implement a way for the api only to be called once when getting parent information. (remove getParentPath)
func (i *Item) ListFiles(client *http.Client) ([]Item, error) {
	url := i.path + ":/children"

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		return nil, err
	}

	var response struct {
		Value []Item `json:"value"`
	}

	json.NewDecoder(resp.Body).Decode(&response)
	for iv := range response.Value {
		response.Value[iv].path = i.path + response.Value[iv].Name
		response.Value[iv].parent = i
	}

	return response.Value, nil
}

func (i *Item) getParentPath(client *http.Client) (string, error) {
	req, _ := http.NewRequest("GET", i.path, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		return "", err
	}

	var response Item

	json.NewDecoder(resp.Body).Decode(&response)

	return response.ParentData.Path, err
}

func Menu2(client *http.Client, sharePoint string) {
	item := NewRootItem(sharePoint)

	for {
		fmt.Print("> ")
		items, err := item.ListFiles(client)

		if err != nil {
			fmt.Println("Error listing files:", err)
			continue
		}

		line := bufio.NewReader(os.Stdin)
		cmd, err := line.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		cmd = strings.TrimSpace(cmd)
		if cmd == "ls" {
			for _, item := range items {
				if item.IsFolder != nil {
					fmt.Println("DIR: ", item.Name)
				} else {
					fmt.Println("FILE: ", item.Name)
				}
			}

			//This whole thing is strange....?
		} else if strings.HasPrefix(cmd, "cd") {
			cmd = strings.TrimPrefix(cmd, "cd")
			cmd = strings.TrimSpace(cmd)

			if cmd != ".." && !directoryExists(items, cmd) {
				fmt.Println("That directory doesn't exist")
				continue
			}
			item = item.ChangeDirectory(cmd)

		} else if strings.HasPrefix(cmd, "pwd") {
			cmd = strings.TrimPrefix(cmd, "pwd")
			cmd = strings.TrimSpace(cmd)

			fmt.Println(item.Pwd())
		}
	}
}
