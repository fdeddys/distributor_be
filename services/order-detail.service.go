package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"

	dto "distribution-system-be/models/dto"

	dbmodels "distribution-system-be/models/dbModels"
)

// OrderDetailService ...
type OrderDetailService struct {
}

// GetDataOrderDetailPage ...
func (o OrderDetailService) GetDataOrderDetailPage(param dto.FilterOrderDetail, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetOrderDetailPage(param, offset, limit)

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
func (o OrderDetailService) Save(orderDetail *dbmodels.SalesOrderDetail) (errCode string, errDesc string) {

	if _, err := database.GetSalesOrderByOrderId(orderDetail.SalesOrderID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SaveSalesOrderDetail(orderDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeleteOrderDetailById ...
func (o OrderDetailService) DeleteOrderDetailByID(orderDetailId int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeleteSalesOrderDetailById(orderDetailId); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateQtyReceive ...
func (o OrderDetailService) UpdateQtyReceive(orderDetailId, qtyReceive int64) (errCode string, errDesc string) {

	if err, errDesc := database.UpdateQtyReceiveSalesOrderDetail(orderDetailId, qtyReceive); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateQtyOrder ...
func (o OrderDetailService) UpdateQtyOrder(orderDetailId, qtyOrder int64) (errCode string, errDesc string) {

	if err, errDesc := database.UpdateQtyOrderSalesOrderDetail(orderDetailId, qtyOrder); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
