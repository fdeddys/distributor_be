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

// GetSalesOrderByOrderId ...
func GetPaymentById(paymentID int64) (dbmodels.Payment, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	payment := dbmodels.Payment{}

	err := db.Preload("Customer").Where(" id = ?  ", paymentID).First(&payment).Error

	return payment, err

}

func GetPaymentOrderBySalesOrderId(salesOrderID int64) (dbmodels.PaymentOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	paymentOrder := dbmodels.PaymentOrder{}

	err := db.Where(" sales_order_id = ?  ", salesOrderID).First(&paymentOrder).Error

	return paymentOrder, err

}

//SavePayment ...
func SavePayment(payment *dbmodels.Payment) (errCode string, errDesc string, id int64) {

	fmt.Println(" Update Payment - ")
	db := GetDbCon()
	db.Debug().LogMode(true)
	payment.PaymentDate = util.GetCurrFormatDateTime()
	r := db.Save(&payment)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		fmt.Println("Error update ", errDesc)
		return
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, payment.ID
}

// GetOrderPage ...
func GetPaymentPage(param dto.FilterPayment, offset, limit int) ([]dbmodels.Payment, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var Payments []dbmodels.Payment
	var total int
	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&Payments).Error
		if err != nil {
			return Payments, 0, err
		}
		return Payments, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysPayments(db, offset, limit, &Payments, param, errQuery)
	go AsyncQueryCountsPayments(db, &total, &Payments, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return Payments, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return Payments, 0, resErrCount
	}
	return Payments, total, nil
}

// AsyncQueryCountsOrders ...
func AsyncQueryCountsPayments(db *gorm.DB, total *int, orders *[]dbmodels.Payment, param dto.FilterPayment, resChan chan error) {

	merchantCode, paymentNo := getParamPayment(param)

	fmt.Println("ISI MERCHANT ", merchantCode, " orderReturnNumber ", paymentNo)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&orders).Where(" COALESCE(payment_no, '') ilike ? AND payment_date between ? and ?  ", paymentNo, param.StartDate, param.EndDate).Count(&*total).Error
	} else {
		err = db.Model(&orders).Where(" COALESCE(payment_no,'') ilike ? ", paymentNo).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrders ...
func AsyncQuerysPayments(db *gorm.DB, offset int, limit int, orders *[]dbmodels.Payment, param dto.FilterPayment, resChan chan error) {

	var err error

	merchantCode, paymentNo := getParamPayment(param)

	fmt.Println("ISI Customer ", merchantCode, " payment no ", paymentNo)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("Customer").Order("payment_date DESC").Offset(offset).Limit(limit).Find(&orders, "  COALESCE(payment_no, '') ilike ? AND payment_date between ? and ?   ", paymentNo, param.StartDate, param.EndDate).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("Customer").Order("payment_date DESC").Find(&orders, "  COALESCE(payment_no,'') ilike ?  ", paymentNo).Error
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

func getParamPayment(param dto.FilterPayment) (merchantCode, paymentNo string) {

	merchantCode = "%"
	paymentNo = "%"

	// paymentNo = param.PaymentNo
	if param.PaymentNo != "" {
		paymentNo = "%" + param.PaymentNo + "%"
	}
	return
}

//RejectPayment ...
func RejectPayment(PaymentId int64) (errCode string, errDesc string) {

	fmt.Println(" Reject Payment by no ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.Payment{}).Where("id =?", PaymentId).Update(dbmodels.Payment{Status: 30})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
