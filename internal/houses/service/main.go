package service

import "homevision/internal/houses/domain"

const (
	PAGES       int    = 10
	IMAGES_PATH string = "./images"
)

var (
	SUPPORTED_IMAGE_FORMATS []string = []string{"jpg", "jpeg", "png"}
)

type HousesService struct {
	housesRepository HousesRepository
}

type HousesRepository interface {
	GetHouses(page int) ([]domain.House, error)
}

func NewHousesService(housesRepository HousesRepository) *HousesService {
	return &HousesService{
		housesRepository: housesRepository,
	}
}
