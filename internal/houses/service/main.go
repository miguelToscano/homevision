package service

import "homevision/internal/houses/domain"

type HousesService struct {
	housesRepository *HousesRepository
}

type HousesRepository interface {
	GetHouses(page int) ([]domain.House, error)
}

func NewHousesService(housesRepository HousesRepository) *HousesService {
	return &HousesService{
		housesRepository: &housesRepository,
	}
}
