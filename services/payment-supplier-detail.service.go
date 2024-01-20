package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// PaymentDetailService ...
type PaymentSupplierDetailService struct {
}

// GetDataPaymentDetailPage ...
func (r PaymentSupplierDetailService) GetDataPaymentDetailPage(param dto.FilterSupplierPayment, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalRec, err := database.GetPaymentSupplierDetailByPaymentID(param.PaymentID, offset, limit)

	if err != nil {
		// res.ErrCode = constants.ERR_CODE_81
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalRec
	res.Page = page
	res.Count = limit

	// res.ErrCode = constants.ERR_CODE_00
	res.Error = constants.ERR_CODE_00_MSG

	return res
}

// SavePaymentDetail ...
func (r PaymentSupplierDetailService) SavePaymentDetail(paymentDetail *dbmodels.PaymentSupplierDetail) (errCode string, errDesc string) {

	if _, err := database.GetPaymentDetailByID(paymentDetail.ID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SavePaymentSupplierDetail(paymentDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeletePaymentDetailByID ...
func (r PaymentSupplierDetailService) DeletePaymentDetailByID(receiveDetailID int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeletePaymentSupplierDetailById(receiveDetailID); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
