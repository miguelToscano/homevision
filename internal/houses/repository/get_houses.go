package repository

import (
	"encoding/json"
	"fmt"
	"homevision/internal/houses/domain"
	"io"
	"net/http"
	"time"
)

const (
	API_URL              = "https://staging.homevision.co/api_project"
	RETRIES              = 5
	TIME_BETWEEN_RETRIES = 1 * time.Second
)

type House struct {
	ID       uint64 `json:"id"`
	Address  string `json:"address"`
	PhotoURL string `json:"photoURL"`
}

type RequestResponse struct {
	Houses []House `json:"houses"`
}

func RetryHTTPRequest(url string, retries int, delay time.Duration) (*http.Response, error) {
	client := &http.Client{}
	var resp *http.Response
	var err error

	for i := 1; i <= retries; i++ {
		resp, err = client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		time.Sleep(delay)
	}

	return nil, fmt.Errorf("failed to fetch the URL after %d attempts: %v", retries, err)
}

func (hs *HousesRepository) GetHouses(page int) ([]domain.House, error) {
	requestResponse, err := RetryHTTPRequest(fmt.Sprintf("%s/houses?page=%d", API_URL, page), RETRIES, TIME_BETWEEN_RETRIES)

	if err != nil {
		return []domain.House{}, err
	}

	defer requestResponse.Body.Close()

	body, err := io.ReadAll(requestResponse.Body)

	if err != nil {
		return []domain.House{}, err
	}

	var unMarshalledRequestResponseBody RequestResponse

	if err := json.Unmarshal(body, &unMarshalledRequestResponseBody); err != nil {
		return []domain.House{}, err
	}

	response := []domain.House{}

	for _, house := range unMarshalledRequestResponseBody.Houses {
		response = append(response, domain.House{
			ID:       house.ID,
			Address:  house.Address,
			PhotoURL: house.PhotoURL,
		})
	}

	return response, nil
}
