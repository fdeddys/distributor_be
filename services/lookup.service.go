package services

import (
	"distribution-system-be/constants"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
)

// LookupService ...
type LookupService struct {
}

// GetLookupByGroup ...
func (h LookupService) GetLookupByGroup(lookupstr string) models.ContentResponse {
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetLookupByGroup(lookupstr)
	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}
	return res
}

// GetPagingLookup ...
func (h LookupService) GetPagingLookup(param dto.FilterLookup, page int, limit int) models.ResponsePagination {

	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetPagingLookup(param, offset, limit)

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

// GetLookupFilter ...
func (h LookupService) GetLookupFilter(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetLookupFilter(id)
	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}
	return res
}

// GetLookupByGroupName ...
func (h LookupService) GetLookupByGroupName(name string) models.ContentResponse {

	fmt.Println("Lookup service ======>")
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetLookupByGroupName(name)
	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}
	return res
}

// SaveLookup ...
func (h LookupService) SaveLookup(lookup *dbmodels.Lookup) models.NoContentResponse {
	// var res models.ResponseSave
	return repository.SaveLookup(*lookup)
}

// GetDistinctLookup ...
func (h LookupService) GetDistinctLookup() models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetDistinctLookup()
	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}
	return res
}

// UpdateLookup ...
func (h LookupService) UpdateLookup(lookup *dbmodels.Lookup) models.NoContentResponse {
	return repository.UpdateLookup(*lookup)
}
