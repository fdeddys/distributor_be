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

// PaymentService ...
type PaymentService struct {
}

// GetDataOrderById ...
func (o PaymentService) GetDataById(paymentID int64) dbmodels.Payment {

	var res dbmodels.Payment
	// var err error
	res, _ = database.GetPaymentById(paymentID)

	return res
}

// GetDataOrderById ...
func (o PaymentService) GetDataPaymentBySalesOrderId(salesOrderID int64) dbmodels.Payment {

	var res dbmodels.Payment

	// var err error
	paymentOrder, err := database.GetPaymentOrderBySalesOrderId(salesOrderID)
	if err != nil {
		return res
	}
	res, _ = database.GetPaymentById(paymentOrder.PaymentID)

	return res
}

// GetDataPage ...
func (o PaymentService) GetDataPage(param dto.FilterPayment, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetPaymentPage(param, offset, limit)

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
func (o PaymentService) Save(payment *dbmodels.Payment) (errCode, errDesc, orderNo string, orderID int64, status int8) {

	if payment.ID == 0 {
		newOrderNo, errCode, errMsg := generateNewPaymentNo(payment.IsCash)
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		payment.PaymentNo = newOrderNo
		payment.Status = constants.STATUS_NEW
	} else {

		// tidak boleh update customer
		curPayment, _ := database.GetPaymentById(payment.ID)
		payment.CustomerID = curPayment.CustomerID
	}
	payment.LastUpdateBy = dto.CurrUser
	payment.LastUpdate = time.Now()

	err, errDesc, newID := database.SavePayment(payment)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, payment.PaymentNo, newID, payment.Status
}

// Approve ...
func (o PaymentService) Approve(payment *dbmodels.Payment) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)

	// orderData := database.GetOrderByOrderNo(payment.)
	// if !(payment.Status == 0 || payment.Status == 10) {
	// 	errCode = constants.ERR_CODE_41
	// 	errDesc = constants.ERR_CODE_41_MSG
	// 	return
	// }

	paymentUpdate, _ := database.GetPaymentById(payment.ID)
	if payment.TotalOrder <= 0 {
		errCode = constants.ERR_CODE_95
		errDesc = constants.ERR_CODE_95_MSG
		return
	}

	if payment.TotalPayment != (payment.TotalOrder - payment.TotalReturn) {
		errCode = constants.ERR_CODE_96
		errDesc = constants.ERR_CODE_96_MSG
		return
	}

	paymentUpdate.Status = constants.STATUS_APPROVE
	paymentUpdate.LastUpdateBy = dto.CurrUser
	paymentUpdate.LastUpdate = time.Now()
	err, errDesc, _ := database.SavePayment(&paymentUpdate)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// Reject ...
func (o PaymentService) Reject(payment *dbmodels.Payment) (errCode, errDesc string) {

	err, errDesc := database.RejectPayment(payment.ID)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func generateNewPaymentNo(isCash bool) (newPaymentNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")

	header := constants.HEADER_PAYMENT_CREDIT
	if isCash {
		header = constants.HEADER_PAYMENT_CASH
	}

	err, number, errdesc := database.AddSequence(bln, thn, header)
	if err != constants.ERR_CODE_00 {
		return "", err, errdesc
	}
	newNumb := fmt.Sprintf("00000%v", number)
	fmt.Println("new numb bef : ", newNumb)
	runes := []rune(newNumb)
	newNumb = string(runes[len(newNumb)-5 : len(newNumb)])
	fmt.Println("new numb after : ", newNumb)

	// newNumb = newNumb[len(newNumb)-5 : len(newNumb)]
	newPaymentNo = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newPaymentNo, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}
