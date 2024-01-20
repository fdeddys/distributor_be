package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// OrderDetailService ...
type ReturnReceiveDetailService struct {
}

// GetReturnOrderDetailPage ...
func (o ReturnReceiveDetailService) GetReturnReceiveDetailPage(param dto.FilterReturnReceive, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReturnReceiveDetailPage(param, offset, limit)

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
func (o ReturnReceiveDetailService) Save(returnReceiveDetail *dbmodels.ReturnReceiveDetail) (errCode string, errDesc string) {

	if _, err := database.GetReceiveByReceiveID(returnReceiveDetail.ReturnReceiveID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SaveReturnReceiveDetail(returnReceiveDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	calculateTotalReturnReceive(returnReceiveDetail.ReturnReceiveID)
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeleteReturnOrderDetailById ...
func (o ReturnReceiveDetailService) DeleteReturnReceiveDetailById(returnReceiveDetailId int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeleteReturnReceiveDetailById(returnReceiveDetailId); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateQtReturn ...
func (o ReturnReceiveDetailService) UpdateQtReturn(returnReceiveDetailId, qty int64) (errCode string, errDesc string) {

	if err, errDesc := database.UpdateQtyReturnReceiveDetail(returnReceiveDetailId, qty); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
