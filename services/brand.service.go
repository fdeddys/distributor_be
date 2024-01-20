package services

import (
	"distribution-system-be/constants"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"time"
)

// BrandService ...
type BrandService struct {
}

// GetBrandFilter ...
func (h BrandService) GetBrandFilter(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetBrandFilter(id)

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

// GetBrandFilterPaging ...
func (h BrandService) GetBrandFilterPaging(param dto.FilterBrand, page int, limit int) models.ResponsePagination {

	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetBrand(param, offset, limit)

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

// GetBrandLike ...
func (h BrandService) GetBrandLike(brandTerms string) models.ContentResponse {
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetBrandLike(brandTerms)

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res
	}

	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	return res
}

// SaveBrand ...
func (h BrandService) SaveBrand(brand *dbmodels.Brand) models.NoContentResponse {
	brand.LastUpdate = time.Now()
	brand.LastUpdateBy = dto.CurrUser
	// var res models.ResponseSave
	res := repository.SaveBrand(*brand)

	return res
}

// UpdateBrand ...
func (h BrandService) UpdateBrand(brand *dbmodels.Brand) models.NoContentResponse {
	brand.LastUpdate = time.Now()
	brand.LastUpdateBy = dto.CurrUser
	// var updatedBrand dbmodels.Brand
	// updatedBrand.ID = brand.ID
	// updatedBrand.Name = brand.Name
	// updatedBrand.Status = brand.Status
	// updatedBrand.LastUpdateBy = brand.LastUpdateBy
	// updatedBrand.LastUpdate = brand.LastUpdate
	// updatedBrand.Code = brand.Code

	res := repository.UpdateBrand(*brand)

	return res
}
