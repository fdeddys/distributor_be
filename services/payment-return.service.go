package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// PaymentReturnService ...
type PaymentReturnService struct {
}

// GetDataPaymentReturnPage ...
func (r PaymentReturnService) GetDataPaymentReturnPage(param dto.FilterPayment) models.ResponsePagination {
	var res models.ResponsePagination

	data, err := database.GetPaymentReturnByPaymentID(param.PaymentID)

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

// SavePaymentReturn ...
func (r PaymentReturnService) SavePaymentReturn(paymentReturn *dbmodels.PaymentReturn) (errCode string, errDesc string) {

	if err, errDesc := database.SavePaymentReturn(paymentReturn); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeletePaymentReturnByID ...
func (r PaymentReturnService) DeletePaymentReturnByID(paymentDetailID int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeletePaymentReturnById(paymentDetailID); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
