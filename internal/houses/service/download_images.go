package service

import (
	"homevision/internal/houses/domain"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	PAGES = 10
)

func (hs *HousesService) DownloadImages() error {

	var wg1 sync.WaitGroup

	for i := 1; i <= PAGES; i++ {
		wg1.Add(1)

		go func(page int) {
			defer wg1.Done()

			houses, _ := (*hs.housesRepository).GetHouses(page)

			var wg2 sync.WaitGroup

			for _, house := range houses {
				wg2.Add(1)

				go func(house domain.House) error {
					defer wg2.Done()

					requestResponse, err := http.Get(house.PhotoURL)

					if err != nil {
						return err
					}

					defer requestResponse.Body.Close()

					file, _ := os.Create(house.GetFileName())

					io.Copy(file, requestResponse.Body)

					return err
				}(house)
			}

			wg2.Wait()
		}(i)
	}

	wg1.Wait()

	return nil
}
