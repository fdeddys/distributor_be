package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// StockMutationDetailService ...
type StockMutationDetailService struct {
}

// GetDataStockMutationDetailPage ...
func (o StockMutationDetailService) GetDataStockMutationDetailPage(param dto.FilterStockMutation, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetStockMutationDetailPage(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = limit

	return res
}

// Save ...
func (o StockMutationDetailService) Save(stockMutationDetail *dbmodels.StockMutationDetail) (errCode string, errDesc string) {

	if _, err := database.GetStockMutationById(stockMutationDetail.MutationID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SaveStockMutationDetail(stockMutationDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeleteStockMutationDetailByID ...
func (o StockMutationDetailService) DeleteStockMutationDetailByID(stockMutationDetailId int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeleteStockMutationDetailById(stockMutationDetailId); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateQty ...
func (o StockMutationDetailService) UpdateQty(stockMutationDetailId, qty int64) (errCode string, errDesc string) {

	if err, errDesc := database.UpdateQtyReceiveSalesOrderDetail(stockMutationDetailId, qty); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
