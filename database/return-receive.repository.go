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

// GetById ...
func GetReturnReceiveById(returnID int64) (dbmodels.ReturnReceive, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	orderReturn := dbmodels.ReturnReceive{}

	err := db.Preload("Supplier").Preload("Warehouse").Preload("Reason").Where(" id = ?  ", returnID).First(&orderReturn).Error

	return orderReturn, err

}

//SaveReceiveReturn ...
func SaveReturnReceive(orderReturn *dbmodels.ReturnReceive) (errCode string, errDesc string, id int64, status int8) {

	fmt.Println("Update Return - ")
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

// GetReturnPage ...
func GetReturnReceivePage(param dto.FilterReturnReceive, offset, limit int) ([]dbmodels.ReturnReceive, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var returnReceives []dbmodels.ReturnReceive
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&returnReceives).Error
		if err != nil {
			return returnReceives, 0, err
		}
		return returnReceives, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysReturnReceives(db, offset, limit, &returnReceives, param, errQuery)
	go AsyncQueryCountsReturnReceives(db, &total, &returnReceives, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return returnReceives, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return returnReceives, 0, resErrCount
	}
	return returnReceives, total, nil
}

// AsyncQueryCountsOrders ...
func AsyncQueryCountsReturnReceives(db *gorm.DB, total *int, orders *[]dbmodels.ReturnReceive, param dto.FilterReturnReceive, resChan chan error) {

	supplierCode, returnNumber := getParamReturnReceive(param)

	fmt.Println("ISI Supplier ", supplierCode, " ReturnNumber ", returnNumber)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&orders).Where(" COALESCE(return_receive_no, '') ilike ? AND return_date between ? and ?  ", returnNumber, param.StartDate, param.EndDate).Count(&*total).Error
	} else {
		err = db.Model(&orders).Where(" COALESCE(return_receive_no,'') ilike ? ", returnNumber).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrders ...
func AsyncQuerysReturnReceives(db *gorm.DB, offset int, limit int, orders *[]dbmodels.ReturnReceive, param dto.FilterReturnReceive, resChan chan error) {

	var err error

	supplierCode, returnNumber := getParamReturnReceive(param)

	fmt.Println("ISI Supplier ", supplierCode, " order no ", returnNumber)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("Supplier").Order("return_date DESC").Offset(offset).Limit(limit).Find(&orders, "  COALESCE(return_receive_no, '') ilike ? AND return_date between ? and ?   ", returnNumber, param.StartDate, param.EndDate).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("Supplier").Order("return_date DESC").Find(&orders, "  COALESCE(return_receive_no,'') ilike ?  ", returnNumber).Error
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

func getParamReturnReceive(param dto.FilterReturnReceive) (supplierCode, returnNumber string) {

	supplierCode = "%"

	returnNumber = param.ReturnNo
	if returnNumber == "" {
		returnNumber = "%"
	} else {
		returnNumber = "%" + param.ReturnNo + "%"
	}

	return
}

//RejectReturnReceive ...
func RejectReturnReceive(salesOrderReturnId int64) (errCode string, errDesc string) {

	fmt.Println(" Reject Return-Sales-Order numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.ReturnReceive{}).Where("id =?", salesOrderReturnId).Update(dbmodels.ReturnReceive{Status: 30})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
