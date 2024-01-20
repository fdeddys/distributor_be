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

type WarehouseService struct {
}

// Get Data Customer Paging
func (m WarehouseService) GetAllWarehouse() models.ContentResponse {
	var res models.ContentResponse

	data, code, msg := database.GetAllWarehouse()
	res.ErrCode = code
	res.ErrDesc = msg
	res.Contents = data
	return res
}

// Get Data Customer Paging
func (m WarehouseService) GetWarehouseByFunc(warehouseIn bool) models.ContentResponse {
	var res models.ContentResponse

	data, code, msg := database.GetWarehouseByFunc(warehouseIn)
	
	res.ErrCode = code
	res.ErrDesc = msg
	res.Contents = data
	return res
}

// GetWarehouseFilter ...
func (s WarehouseService) GetWarehouseFilter(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetWarehouseFilter(id)

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

// GetWarehouseFilterPaging ...
func (s WarehouseService) GetWarehouseFilterPaging(param dto.FilterPaging, page int, limit int) models.ResponsePagination {

	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetWarehouse(param, offset, limit)

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

// GetWarehouseLike ...
func (s WarehouseService) GetWarehouseLike(terms string) models.ContentResponse {
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetWarehouseLike(terms)

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
func (s WarehouseService) SaveWarehouse(warehouse *dbmodels.Warehouse) models.NoContentResponse {
	warehouse.LastUpdate = time.Now()
	warehouse.LastUpdateBy = dto.CurrUser
	return repository.SaveWarehouse(*warehouse)
}

// UpdateBrand ...
func (s WarehouseService) UpdateWarehouse(warehouse *dbmodels.Warehouse) models.NoContentResponse {
	warehouse.LastUpdate = time.Now()
	warehouse.LastUpdateBy = dto.CurrUser

	return repository.UpdateWarehouse(*warehouse)

}
