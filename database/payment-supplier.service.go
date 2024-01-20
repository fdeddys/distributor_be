package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetOrderPage ...
func GetPaymentSupplierPage(param dto.FilterSupplierPayment, offset, limit int) ([]dbmodels.PaymentSupplier, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var paymentSuppliers []dbmodels.PaymentSupplier
	var total int
	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&paymentSuppliers).Error
		if err != nil {
			return paymentSuppliers, 0, err
		}
		return paymentSuppliers, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysPaymentSuppliers(db, offset, limit, &paymentSuppliers, param, errQuery)
	go AsyncQueryCountsPaymentSuppliers(db, &total, &paymentSuppliers, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return paymentSuppliers, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return paymentSuppliers, 0, resErrCount
	}
	return paymentSuppliers, total, nil
}

// AsyncQueryCountsOrders ...
func AsyncQueryCountsPaymentSuppliers(db *gorm.DB, total *int, orders *[]dbmodels.PaymentSupplier, param dto.FilterSupplierPayment, resChan chan error) {

	_, paymentNo, byStatus := getParamPaymentSupplier(param)

	fmt.Println("  Payment no ", paymentNo)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&orders).Where(" COALESCE(payment_no, '') ilike ? AND ( payment_date between ? and ? ) AND  ( (status = ?) or ( not ?) ) ", paymentNo, param.StartDate, param.EndDate, param.PaymentStatus, byStatus).Count(&*total).Error
	} else {
		err = db.Model(&orders).Where(" COALESCE(payment_no,'') ilike ?  AND  ( (status = ?) or ( not ?) )", paymentNo, param.PaymentStatus, byStatus).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrders ...
func AsyncQuerysPaymentSuppliers(db *gorm.DB, offset int, limit int, orders *[]dbmodels.PaymentSupplier, param dto.FilterSupplierPayment, resChan chan error) {

	var err error

	supplierCode, paymentNo, byStatus := getParamPaymentSupplier(param)

	fmt.Println("ISI Customer ", supplierCode, " payment no ", paymentNo)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("Supplier").Preload("Supplier.Bank").Order("payment_date DESC").Offset(offset).Limit(limit).Find(&orders, "  COALESCE(payment_no, '') ilike ? AND ( payment_date between ? and ? ) AND  ( (status = ?) or ( not ?) ) ", paymentNo, param.StartDate, param.EndDate, param.PaymentStatus, byStatus).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("Supplier").Order("payment_date DESC").Find(&orders, "  COALESCE(payment_no,'') ilike ?  AND  ( (status = ?) or ( not ?) ) ", paymentNo, param.PaymentStatus, byStatus).Error
		if err != nil {
			fmt.Println("error --> ", err)
		}
		fmt.Println("order--> ", orders)
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

func getParamPaymentSupplier(param dto.FilterSupplierPayment) (supplierCode, paymentNo string, byStatus bool) {

	supplierCode = "%"
	paymentNo = "%"

	byStatus = true
	if param.PaymentStatus == 0 {
		byStatus = false
	}

	// paymentNo = param.PaymentNo
	if param.PaymentNo != "" {
		paymentNo = "%" + param.PaymentNo + "%"
	}
	return
}

// GetSalesOrderByOrderId ...
func GetPaymentSupplierById(paymentID int64) (dbmodels.PaymentSupplier, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	paymentSupplier := dbmodels.PaymentSupplier{}

	err := db.Preload("Supplier").Where(" id = ?  ", paymentID).First(&paymentSupplier).Error

	return paymentSupplier, err

}

//SavePayment ...
func SavePaymentSupplier(paymentSupplier *dbmodels.PaymentSupplier) (errCode string, errDesc string, id int64) {

	fmt.Println(" Update Payment - ")
	db := GetDbCon()
	db.Debug().LogMode(true)
	paymentSupplier.PaymentDate = util.GetCurrFormatDateTime()
	r := db.Save(&paymentSupplier)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		fmt.Println("Error update ", errDesc)
		return
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, paymentSupplier.ID
}
