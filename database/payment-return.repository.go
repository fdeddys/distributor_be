package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
)

//SavePaymentDetailNo ...
func SavePaymentReturn(paymentReturn *dbmodels.PaymentReturn) (errCode string, errDesc string) {

	fmt.Println(" Add Payment order ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	paymentReturn.LastUpdateBy = dto.CurrUser
	paymentReturn.LastUpdate = util.GetCurrDate()
	db.Save(&paymentReturn)

	db.Model(&dbmodels.ReturnSalesOrder{}).
		Where("id = ?", paymentReturn.ReturnSalesOrderID).
		Update(dbmodels.ReturnSalesOrder{
			IsPaid: true,
		})

	payment := dbmodels.Payment{}
	db.Where("id = ?", paymentReturn.PaymentID).Find(&payment)
	payment.TotalReturn = payment.TotalReturn + paymentReturn.Total
	payment.LastUpdateBy = dto.CurrUser
	payment.LastUpdate = util.GetCurrDate()
	db.Save(&payment)

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func DeletePaymentReturnById(paymentDetailID int64) (errCode string, errDesc string) {

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
	paymentReturn := dbmodels.PaymentReturn{}
	tx.Where(" id = ?  ", paymentDetailID).First(&paymentReturn)

	tx.Where("id = ? ", paymentDetailID).Delete(dbmodels.PaymentReturn{})

	db.Model(&dbmodels.ReturnSalesOrder{}).Where("id = ?", paymentReturn.ReturnSalesOrderID).Update("is_paid", false)

	// tx.Model(&dbmodels.ReturnSalesOrder{}).
	// 	Where("id = ?", paymentReturn.ReturnSalesOrderID).
	// 	Update(dbmodels.SalesOrder{
	// 		IsPaid: false,
	// 	})

	payment := dbmodels.Payment{}
	tx.Where("id = ?", paymentReturn.PaymentID).Find(&payment)
	payment.TotalOrder = payment.TotalOrder - paymentReturn.Total
	payment.LastUpdateBy = dto.CurrUser
	payment.LastUpdate = util.GetCurrDate()
	tx.Save(&payment)

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// GetPaymentDetailByPaymentNo ...
func GetPaymentReturnByPaymentID(paymentID int64) ([]dbmodels.PaymentReturn, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	var paymentDetail []dbmodels.PaymentReturn

	err := db.Preload("ReturnSalesOrder").Where(" payment_id = ?  ", paymentID).Find(&paymentDetail).Error

	return paymentDetail, err

}
