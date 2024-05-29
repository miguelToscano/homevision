package service

import (
	"fmt"
	"os"
	"strings"
)

func (hs *HousesService) GetDownloadedImages() []string {

	photosURLs := []string{}

	dir := "."

	// Open the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return []string{}
	}

	// Loop through directory entries
	for _, file := range files {
		// Check if the entry is a file and has a .png extension
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			// Print the PNG file name
			photosURLs = append(photosURLs, file.Name())
		}
	}

	return photosURLs
}
