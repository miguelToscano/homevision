package handler

type HousesHandler struct {
	housesService *HousesService
}

type HousesService interface {
	DownloadImages() error
	GetDownloadedImages() []string
}

func NewHousesHandler(housesService HousesService) *HousesHandler {
	return &HousesHandler{
		housesService: &housesService,
	}
}
