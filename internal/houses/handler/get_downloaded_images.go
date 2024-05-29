package handler

func (hh *HousesHandler) GetDownloadedImages() []string {
	downloadedImages := (*hh.housesService).GetDownloadedImages()
	return downloadedImages
}
