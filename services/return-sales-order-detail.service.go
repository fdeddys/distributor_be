package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// OrderDetailService ...
type ReturnOrderDetailService struct {
}

// GetReturnOrderDetailPage ...
func (o ReturnOrderDetailService) GetReturnOrderDetailPage(param dto.FilterOrderReturnDetail, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReturnOrderDetailPage(param, offset, limit)

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
func (o ReturnOrderDetailService) Save(returnOrderDetail *dbmodels.ReturnSalesOrderDetail) (errCode string, errDesc string) {

	if _, err := database.GetSalesOrderByOrderId(returnOrderDetail.ReturnSalesOrderID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SaveReturnOrderDetail(returnOrderDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeleteReturnOrderDetailById ...
func (o ReturnOrderDetailService) DeleteReturnOrderDetailById(orderDetailId int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeleteReturnOrderDetailById(orderDetailId); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateQtReturn ...
func (o ReturnOrderDetailService) UpdateQtReturn(retrunOrderDetailId, qtyReceive int64) (errCode string, errDesc string) {

	if err, errDesc := database.UpdateQtyReturnSalesOrderDetail(retrunOrderDetailId, qtyReceive); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
