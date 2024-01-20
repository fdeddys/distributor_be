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

// GetAllReceiveReturnDetail ...
func GetAllReceiveReturnDetail(orderReturnID int64) []dbmodels.ReturnReceiveDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var orderReturnDetails []dbmodels.ReturnReceiveDetail

	db.Preload("Product").Preload("UOM").Find(&orderReturnDetails, " return_receive_id = ?  ", orderReturnID)

	return orderReturnDetails
}

// GetReturnOrderDetailPage ...
func GetReturnReceiveDetailPage(param dto.FilterReturnReceive, offset, limit int) ([]dbmodels.ReturnReceiveDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var returnReceiveDetails []dbmodels.ReturnReceiveDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&returnReceiveDetails).Error
		if err != nil {
			return returnReceiveDetails, 0, err
		}
		return returnReceiveDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysReceiveReturnDetails(db, offset, limit, &returnReceiveDetails, param.ReturnID, errQuery)
	go AsyncQueryCountsReceiveReturnDetails(db, &total, param.ReturnID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return returnReceiveDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return returnReceiveDetails, 0, resErrCount
	}
	return returnReceiveDetails, total, nil
}

// AsyncQueryCountsOrderDetails ...
func AsyncQueryCountsReceiveReturnDetails(db *gorm.DB, total *int, returnReceiveID int64, offset int, limit int, resChan chan error) {

	var err error

	err = db.Model(&dbmodels.ReturnReceiveDetail{}).Where("return_receive_id = ?", returnReceiveID).Count(total).Error
	// Offset(offset).
	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrderReturnDetails ...
func AsyncQuerysReceiveReturnDetails(db *gorm.DB, offset int, limit int, returnOrderDetails *[]dbmodels.ReturnReceiveDetail, returnReceiveID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).Preload("Product").Preload("UOM").Order("id desc").Find(&returnOrderDetails, " return_receive_id = ? ", returnReceiveID).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("data--> ", returnOrderDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//SaveReturnReceiveDetail ...
func SaveReturnReceiveDetail(returnOrderDetail *dbmodels.ReturnReceiveDetail) (errCode string, errDesc string) {

	fmt.Println(" Update Return Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Save(&returnOrderDetail); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", returnOrderDetail.ID)
	return

}

// DeleteReturnReceiveDetailById ...
func DeleteReturnReceiveDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Return Detail  ------------------------------------------ ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Where("id = ? ", id).Delete(dbmodels.ReturnReceiveDetail{}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = r.Error.Error()
	}

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", id)
	return

}

//UpdateQtyReturnReceiveDetail ...
func UpdateQtyReturnReceiveDetail(returnSaleOrderDetailId int64, qty int64) (errCode string, errDesc string) {

	fmt.Println(" Update Qty return Receive Detail  -- ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("id = %v, qty = %v", returnSaleOrderDetailId, qty)

	if r := db.Model(&dbmodels.ReturnReceiveDetail{}).
		Where("id = ?", returnSaleOrderDetailId).
		Update(dbmodels.ReturnReceiveDetail{
			Qty:          qty,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}
