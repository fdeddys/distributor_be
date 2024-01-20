package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// PaymentOrderService ...
type PaymentOrderService struct {
}

// GetDataPaymentOrderPage ...
func (r PaymentOrderService) GetDataPaymentOrderPage(param dto.FilterPayment) models.ResponsePagination {
	var res models.ResponsePagination

	data, err := database.GetPaymentOrderByPaymentID(param.PaymentID)

	if err != nil {
		// res.ErrCode = constants.ERR_CODE_81
		res.Error = err.Error()
		res.Contents = data
		return res
	}

	res.Contents = data
	// res. = constants.ERR_CODE_00
	res.Error = constants.ERR_CODE_00_MSG

	return res
}

// SavePaymentOrder ...
func (r PaymentOrderService) SavePaymentOrder(paymentOrder *dbmodels.PaymentOrder) (errCode string, errDesc string) {

	// db := database.GetDbCon()
	// tx := db.Begin()

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	// tx.Model(&dbmodels.SalesOrder{}).
	// 	Where("id = ?", paymentOrder.SalesOrderID).
	// 	Update(dbmodels.SalesOrder{
	// 		IsPaid: true,
	// 	})

	// tx.Save(paymentOrder)

	// tx.Commit()

	order, _ := database.GetSalesOrderByOrderId(paymentOrder.SalesOrderID)
	if !(order.Status == 0 || order.Status == 10 || order.Status == 20) {
		return constants.ERR_CODE_41, constants.ERR_CODE_41_MSG
	}

	if err, errDesc := database.SavePaymentOrder(paymentOrder); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// DeletePaymentOrderByID ...
func (r PaymentOrderService) DeletePaymentOrderByID(paymentDetailID int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeletePaymentOrderById(paymentDetailID); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
