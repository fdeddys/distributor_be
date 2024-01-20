package services

import (
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
)

// ProductGroupService ...
type ProductGroupService struct {
}

// GetProductGroupPaging ...
func (h ProductGroupService) GetProductGroupPaging(param dto.FilterProductGroup, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetProductGroupPaging(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = len(data)

	return res
}

// GetProductGroupDetails ...
func (h ProductGroupService) GetProductGroupDetails(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetProductGroupDetails(id)

	if err != nil {
		res.Contents = nil
		res.ErrCode = "02"
		res.ErrDesc = "Error query data to DB"
		return res
	}

	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	return res
}

// SaveProductGroup ...
func (h ProductGroupService) SaveProductGroup(productGroup *dbmodels.ProductGroup) models.NoContentResponse {
	// productGroup.LastUpdate = time.Now()
	// productGroup.LastUpdateBy = dto.CurrUser

	// var res models.ResponseSave
	res := repository.SaveProductGroup(*productGroup)

	return res
}

// UpdateProductGroup ...
func (h ProductGroupService) UpdateProductGroup(productGroup *dbmodels.ProductGroup) models.NoContentResponse {
	// productGroup.LastUpdate = time.Now()
	// productGroup.LastUpdateBy = dto.CurrUser

	res := repository.UpdateProductGroup(*productGroup)

	return res
}
