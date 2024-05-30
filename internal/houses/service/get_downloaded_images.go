package service

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func (hs *HousesService) GetDownloadedImages() []string {

	images := []string{}

	files, err := os.ReadDir(IMAGES_PATH)

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return []string{}
	}

	for _, file := range files {

		if !file.IsDir() && slices.Contains(SUPPORTED_IMAGE_FORMATS, strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1]) {
			images = append(images, file.Name())
		}
	}

	slices.SortFunc(images, func(i, j string) int {
		first, err := strconv.ParseInt(strings.Split(i, "-")[0], 10, 64)

		if err != nil {
			fmt.Println("Error parsing int:", err)
			return 0
		}

		second, err := strconv.ParseInt(strings.Split(j, "-")[0], 10, 64)

		if err != nil {
			fmt.Println("Error parsing int:", err)
			return 0
		}

		return cmp.Compare(first, second)
	})

	return images
}
