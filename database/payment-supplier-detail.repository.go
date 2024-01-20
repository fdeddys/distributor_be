package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetPaymentDetailByPaymentNo ...
func GetPaymentSupplierDetailByPaymentID(paymentID int64, offset, limit int) ([]dbmodels.PaymentSupplierDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	var paymentSupplierDetails []dbmodels.PaymentSupplierDetail

	var total int
	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&paymentSupplierDetails).Error
		if err != nil {
			return paymentSupplierDetails, 0, err
		}
		return paymentSupplierDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysPaymentSupplierDetails(db, offset, limit, &paymentSupplierDetails, paymentID, errQuery)
	go AsyncQueryCountsPaymentSupplierDetails(db, &total, &paymentSupplierDetails, paymentID, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return paymentSupplierDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return paymentSupplierDetails, 0, resErrCount
	}
	return paymentSupplierDetails, total, nil

}

// AsyncQueryCountsOrders ...
func AsyncQueryCountsPaymentSupplierDetails(db *gorm.DB, total *int, orders *[]dbmodels.PaymentSupplierDetail, paymentID int64, resChan chan error) {

	var err error

	err = db.Model(&orders).Where(" payment_supplier_id =?", paymentID).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrders ...
func AsyncQuerysPaymentSupplierDetails(db *gorm.DB, offset int, limit int, orders *[]dbmodels.PaymentSupplierDetail, paymentID int64, resChan chan error) {

	var err error

	err = db.Preload("Receive").Offset(offset).Limit(limit).Find(&orders, "  payment_supplier_id =?  ", paymentID).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//SavePaymentDetail ...
func SavePaymentSupplierDetail(paymentDetail *dbmodels.PaymentSupplierDetail) (errCode string, errDesc string) {

	fmt.Println(" Add Payment Order Detail ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	paymentDetail.LastUpdateBy = dto.CurrUser
	paymentDetail.LastUpdate = util.GetCurrDate()
	db.Save(&paymentDetail)

	db.Model(&dbmodels.Receive{}).
		Where("id = ?", paymentDetail.ReceiveID).
		Update(dbmodels.Receive{
			IsPaid: true,
			Status: constants.STATUS_PAID,
		})

	paymentSupplier := dbmodels.PaymentSupplier{}
	db.Where("id = ?", paymentDetail.PaymentSupplierID).Find(&paymentSupplier)
	paymentSupplier.Total = paymentSupplier.Total + paymentDetail.Total
	paymentSupplier.LastUpdateBy = dto.CurrUser
	paymentSupplier.LastUpdate = util.GetCurrDate()
	db.Save(&paymentSupplier)

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

func DeletePaymentSupplierDetailById(paymentDetailID int64) (errCode string, errDesc string) {

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
	paymentDetail := dbmodels.PaymentSupplierDetail{}
	tx.Where(" id = ?  ", paymentDetailID).First(&paymentDetail)

	// cancel is paid
	db.Model(&dbmodels.Receive{}).Where("id = ?", paymentDetail.ReceiveID).Update(dbmodels.Receive{
		IsPaid: false,
		Status: constants.STATUS_APPROVE,
	})

	// DELETE Payment order - DETAIL
	tx.Where("id = ? ", paymentDetailID).Delete(dbmodels.PaymentSupplierDetail{})

	// update total money
	paymentSupplier := dbmodels.PaymentSupplier{}
	tx.Where("id = ?", paymentDetail.PaymentSupplierID).Find(&paymentSupplier)

	paymentSupplier.Total = paymentSupplier.Total - paymentDetail.Total
	paymentSupplier.LastUpdateBy = dto.CurrUser
	paymentSupplier.LastUpdate = util.GetCurrDate()
	tx.Save(&paymentSupplier)

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

func TotalPaymentDetailByPaymentID(paymentSupplierID int64) float32 {
	total := float32(0)

	var orderDetails []dbmodels.PaymentSupplierDetail
	db := GetDbCon()
	db.Debug().LogMode(true)

	err := db.Find(&orderDetails, "  payment_supplier_id =?  ", paymentSupplierID).Error
	if err != nil {
		return total
	}
	if len(orderDetails) < 1 {
		return total
	}

	for _, orderDetail := range orderDetails {
		total1 := orderDetail.Total
		total += total1
	}

	return total
}
