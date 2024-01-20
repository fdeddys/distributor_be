package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// PaymentDetailService ...
type PaymentDetailService struct {
}

// GetDataPaymentDetailPage ...
func (r PaymentDetailService) GetDataPaymentDetailPage(param dto.FilterPayment) models.ResponsePagination {
	var res models.ResponsePagination

	data, err := database.GetPaymentDetailByPaymentID(param.PaymentID)

	if err != nil {
		// res.ErrCode = constants.ERR_CODE_81
		res.Error = err.Error()
		return res
	}

	res.Contents = data

	// res.ErrCode = constants.ERR_CODE_00
	res.Error = constants.ERR_CODE_00_MSG

	return res
}

// SavePaymentDetail ...
func (r PaymentDetailService) SavePaymentDetail(paymentDetail *dbmodels.PaymentDetail) (errCode string, errDesc string) {

	if _, err := database.GetPaymentDetailByID(paymentDetail.ID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SavePaymentDetail(paymentDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeletePaymentDetailByID ...
func (r PaymentDetailService) DeletePaymentDetailByID(receiveDetailID int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeletePaymentDetailById(receiveDetailID); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
