package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetSalesOrderByOrderId ...
func GetSalesOrderReturnById(orderReturnID int64) (dbmodels.ReturnSalesOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	orderReturn := dbmodels.ReturnSalesOrder{}

	err := db.Preload("Customer").Preload("Warehouse").Preload("Reason").Preload("Salesman").Where(" id = ?  ", orderReturnID).First(&orderReturn).Error

	return orderReturn, err

}

//SaveSalesOrderReturn ...
func SaveSalesOrderReturn(orderReturn *dbmodels.ReturnSalesOrder) (errCode string, errDesc string, id int64, status int8) {

	fmt.Println(" Update Sales Order Return - ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Save(&orderReturn)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, orderReturn.ID, orderReturn.Status
}

// GetOrderPage ...
func GetReturnSalesOrderPage(param dto.FilterOrderReturnDetail, offset, limit int) ([]dbmodels.ReturnSalesOrder, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var returnSalesOrders []dbmodels.ReturnSalesOrder
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&returnSalesOrders).Error
		if err != nil {
			return returnSalesOrders, 0, err
		}
		return returnSalesOrders, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysReturnSalesOrders(db, offset, limit, &returnSalesOrders, param, errQuery)
	go AsyncQueryCountsReturnSalesOrders(db, &total, &returnSalesOrders, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return returnSalesOrders, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return returnSalesOrders, 0, resErrCount
	}
	return returnSalesOrders, total, nil
}

// AsyncQueryCountsOrders ...
func AsyncQueryCountsReturnSalesOrders(db *gorm.DB, total *int, orders *[]dbmodels.ReturnSalesOrder, param dto.FilterOrderReturnDetail, resChan chan error) {

	merchantCode, orderNumber := getParamReturnSalesOrder(param)

	fmt.Println("ISI MERCHANT ", merchantCode, " orderReturnNumber ", orderNumber)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&orders).Where(" COALESCE(return_no, '') ilike ? AND return_date between ? and ?  ", orderNumber, param.StartDate, param.EndDate).Count(&*total).Error
	} else {
		err = db.Model(&orders).Where(" COALESCE(return_no,'') ilike ? ", orderNumber).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrders ...
func AsyncQuerysReturnSalesOrders(db *gorm.DB, offset int, limit int, orders *[]dbmodels.ReturnSalesOrder, param dto.FilterOrderReturnDetail, resChan chan error) {

	var err error

	merchantCode, orderNumber := getParamReturnSalesOrder(param)

	fmt.Println("ISI MERCHANT ", merchantCode, " order no ", orderNumber)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("Customer").Preload("Salesman").Order("return_date DESC, id desc").Offset(offset).Limit(limit).Find(&orders, "  COALESCE(return_no, '') ilike ? AND return_date between ? and ?   ", orderNumber, param.StartDate, param.EndDate).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("Customer").Preload("Salesman").Order("return_date DESC, id desc").Find(&orders, "  COALESCE(return_no,'') ilike ?  ", orderNumber).Error
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

func getParamReturnSalesOrder(param dto.FilterOrderReturnDetail) (merchantCode, orderNumber string) {

	merchantCode = "%"

	orderNumber = param.OrderReturnNo
	if orderNumber == "" {
		orderNumber = "%"
	} else {
		orderNumber = "%" + param.OrderReturnNo + "%"
	}

	return
}

//RejectReturnSalesOrder ...
func RejectReturnSalesOrder(salesOrderReturnId int64) (errCode string, errDesc string) {

	fmt.Println(" Reject Return-Sales-Order numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.ReturnSalesOrder{}).Where("id =?", salesOrderReturnId).Update(dbmodels.ReturnSalesOrder{Status: 30})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func GetSalesOrderReturnForPayment(param dto.FilterOrderReturnDetail, offset, limit int) ([]dbmodels.ReturnSalesOrder, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var returns []dbmodels.ReturnSalesOrder
	var total int

	total = 0
	err := db.Offset(offset).Limit(limit).Preload("Customer").Order("return_date DESC").Find(&returns, "  ( status  in ('20','40') ) AND (is_paid is null or is_paid = false) AND customer_id = ? ", param.CustomerID).Error

	if err != nil {
		return returns, total, err
	}

	var cekReturn []dbmodels.ReturnSalesOrder
	db.Model(&cekReturn).Where("  ( status  in  ('20','40')  ) AND (is_paid is null or is_paid = false) AND customer_id = ? ", param.CustomerID).Count(&total)

	return returns, total, nil
}
