package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"time"
)

// PaymentSupplierService ...
type PaymentSupplierService struct {
}

// GetDataPage ...
func (o PaymentSupplierService) GetDataPage(param dto.FilterSupplierPayment, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetPaymentSupplierPage(param, offset, limit)

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

// GetDataOrderById ...
func (o PaymentSupplierService) GetDataById(paymentID int64) dbmodels.PaymentSupplier {

	var res dbmodels.PaymentSupplier
	// var err error
	res, _ = database.GetPaymentSupplierById(paymentID)

	return res
}

// Save ...
func (o PaymentSupplierService) Save(paymentSupplier *dbmodels.PaymentSupplier) (errCode, errDesc, orderNo string, orderID int64, status int8) {

	if paymentSupplier.ID == 0 {
		newOrderNo, errCode, errMsg := generateNewPaymentSupplierNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		paymentSupplier.PaymentNo = newOrderNo
		paymentSupplier.Status = constants.STATUS_NEW
	}
	// else {
	// 	// tidak boleh update customer
	// 	curPayment, _ := database.GetPaymentSupplierById(paymentSupplier.ID)
	// 	paymentSupplier.SupplierID = curPayment.SupplierID
	// }
	paymentSupplier.Total = database.TotalPaymentDetailByPaymentID(paymentSupplier.ID)
	paymentSupplier.LastUpdateBy = dto.CurrUser
	paymentSupplier.LastUpdate = time.Now()
	err, errDesc, newID := database.SavePaymentSupplier(paymentSupplier)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, paymentSupplier.PaymentNo, newID, paymentSupplier.Status
}

func generateNewPaymentSupplierNo() (newPaymentNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")

	header := constants.HEADER_PAYMENT_SUPPLIER

	err, number, errdesc := database.AddSequence(bln, thn, header)
	if err != constants.ERR_CODE_00 {
		return "", err, errdesc
	}
	newNumb := fmt.Sprintf("00000%v", number)
	fmt.Println("new numb bef : ", newNumb)
	runes := []rune(newNumb)
	newNumb = string(runes[len(newNumb)-5 : len(newNumb)])
	fmt.Println("new numb after : ", newNumb)

	newPaymentNo = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newPaymentNo, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// Approve ...
func (o PaymentSupplierService) Approve(paymentSupplier *dbmodels.PaymentSupplier) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)

	paymentUpdate, _ := database.GetPaymentSupplierById(paymentSupplier.ID)

	// if paymentSupplier.Total <= 0 {
	// 	errCode = constants.ERR_CODE_96
	// 	errDesc = constants.ERR_CODE_96_MSG
	// 	return
	// }

	paymentUpdate.Status = constants.STATUS_APPROVE
	paymentUpdate.LastUpdateBy = dto.CurrUser
	paymentUpdate.LastUpdate = time.Now()
	err, errDesc, _ := database.SavePaymentSupplier(&paymentUpdate)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
