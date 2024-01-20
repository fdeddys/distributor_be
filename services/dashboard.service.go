package services

import (
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	"distribution-system-be/models/dto"
)

// DashboardService ...
type DashboardService struct {
}

// GetQtyOrder ...
func (d DashboardService) GetQtyOrder(param dto.FilterDto) models.ContentResponse {
	var res models.ContentResponse

	resContent, _ := repository.GetQtyOrd(param.StartDate, param.EndDate)

	res.Contents = resContent
	return res
}
