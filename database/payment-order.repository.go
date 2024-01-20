package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
)

//SavePaymentDetailNo ...
func SavePaymentOrder(paymentOrder *dbmodels.PaymentOrder) (errCode string, errDesc string) {

	fmt.Println(" Add Payment order ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	paymentOrder.LastUpdateBy = dto.CurrUser
	paymentOrder.LastUpdate = util.GetCurrDate()
	db.Save(&paymentOrder)

	db.Model(&dbmodels.SalesOrder{}).
		Where("id = ?", paymentOrder.SalesOrderID).
		Update(dbmodels.SalesOrder{
			IsPaid: true,
			Status: constants.STATUS_PAID,
		})

	payment := dbmodels.Payment{}
	db.Where("id = ?", paymentOrder.PaymentID).Find(&payment)
	payment.TotalOrder = payment.TotalOrder + paymentOrder.Total
	payment.LastUpdateBy = dto.CurrUser
	payment.LastUpdate = util.GetCurrDate()
	db.Save(&payment)

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func DeletePaymentOrderById(paymentDetailID int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Payment Detail  ------------------------------------------  ", paymentDetailID)

	db := GetDbCon()
	db.Debug().LogMode(true)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// FIND PAYMENT DETAIL
	paymentOrder := dbmodels.PaymentOrder{}
	tx.Where(" id = ?  ", paymentDetailID).First(&paymentOrder)

	// DELETE Payment order - DETAIL
	tx.Where("id = ? ", paymentDetailID).Delete(dbmodels.PaymentOrder{})

	// cancel is paid
	db.Model(&dbmodels.SalesOrder{}).Where("id = ?", paymentOrder.SalesOrderID).Update("is_paid", false)

	// tx.Model(&dbmodels.SalesOrder{}).
	// 	Where("id = ?", paymentOrder.SalesOrderID).
	// 	Update(dbmodels.SalesOrder{
	// 		IsPaid: false,
	// 	})

	// update total money
	payment := dbmodels.Payment{}
	tx.Where("id = ?", paymentOrder.PaymentID).Find(&payment)

	payment.TotalOrder = payment.TotalOrder - paymentOrder.Total
	payment.LastUpdateBy = dto.CurrUser
	payment.LastUpdate = util.GetCurrDate()
	tx.Save(&payment)

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// GetPaymentDetailByPaymentNo ...
func GetPaymentOrderByPaymentID(paymentID int64) ([]dbmodels.PaymentOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	var paymentOrder []dbmodels.PaymentOrder

	err := db.Preload("SalesOrder").Where(" payment_id = ?  ", paymentID).Find(&paymentOrder).Error

	return paymentOrder, err

}
