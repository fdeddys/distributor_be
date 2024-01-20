package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"time"
)

type SalesmanService struct {
}

// Get Data Customer Paging
func (s SalesmanService) GetAllSalesman() models.ContentResponse {
	var res models.ContentResponse

	data, code, msg := database.GetAllSalesman()
	res.ErrCode = code
	res.ErrDesc = msg
	res.Contents = data
	return res
}

// GetSalesmanFilter ...
func (s SalesmanService) GetSalesmanFilter(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetSalesmanFilter(id)

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

// GetSalesmanFilterPaging ...
func (s SalesmanService) GetSalesmanFilterPaging(param dto.FilterPaging, page int, limit int) models.ResponsePagination {

	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetSalesman(param, offset, limit)

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

// GetSalesmanLike ...
func (s SalesmanService) GetSalesmanLike(terms string) models.ContentResponse {
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetSalesmanLike(terms)

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
func (s SalesmanService) SaveSalesman(salesman *dbmodels.Salesman) models.NoContentResponse {
	salesman.LastUpdate = time.Now()
	salesman.LastUpdateBy = dto.CurrUser
	return repository.SaveSalesman(*salesman)
}

// UpdateBrand ...
func (s SalesmanService) UpdateSalesman(salesman *dbmodels.Salesman) models.NoContentResponse {
	salesman.LastUpdate = time.Now()
	salesman.LastUpdateBy = dto.CurrUser

	return repository.UpdateSalesman(*salesman)

}
