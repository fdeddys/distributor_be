package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	"fmt"

	"distribution-system-be/models/dto"
)

// DirectPaymentService ...
type DirectPaymentService struct {
}

// GetDataPage ...
func (o DirectPaymentService) GetDataPage(param dto.FilterPaymentDirect, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetPaymentDirectPage(param, offset, limit)

	fmt.Println("Hasil direct =>", data)
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

// Approve ...
func (o DirectPaymentService) Approve(paymentDirect *dto.PaymentDirectModel) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)

	errCode, errDesc = database.ApprovePaymentDirect(paymentDirect)
	if errCode != constants.ERR_CODE_00 {
		return errCode, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// Reject ...
func (o DirectPaymentService) Reject(paymentDirect *dto.PaymentDirectModel) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)

	errCode, errDesc = database.RejectPaymentDirect(paymentDirect)
	if errCode != constants.ERR_CODE_00 {
		return errCode, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
