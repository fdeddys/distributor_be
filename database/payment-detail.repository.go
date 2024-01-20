package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
)

//SavePaymentDetail ...
func SavePaymentDetail(paymentDetail *dbmodels.PaymentDetail) (errCode string, errDesc string) {

	fmt.Println(" Add Payment detail ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	paymentDetail.LastUpdateBy = dto.CurrUser
	paymentDetail.LastUpdate = util.GetCurrDate()
	r := db.Save(&paymentDetail)
	if r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func DeletePaymentDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Payment Detail  ------------------------------------------  ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", id)

	if r := db.Where("id = ? ", id).Delete(dbmodels.PaymentDetail{}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}

// GetPaymentDetailByPaymentNo ...
func GetPaymentDetailByPaymentID(paymentID int64) ([]dbmodels.PaymentDetail, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	var paymentDetail []dbmodels.PaymentDetail

	err := db.Preload("PaymentType").Where(" payment_id = ?  ", paymentID).Find(&paymentDetail).Error

	return paymentDetail, err

}

func GetPaymentDetailByID(id int64) (dbmodels.PaymentDetail, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	paymentDetail := dbmodels.PaymentDetail{}

	err := db.Where(" id = ?  ", id).Error

	return paymentDetail, err

}
