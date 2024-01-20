package database

import (
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetAllDataDetailAdjustment ...
func GetAllDataDetailAdjustment(adjustmentID int64) []dbmodels.AdjustmentDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var adjustmentDetails []dbmodels.AdjustmentDetail

	db.Preload("Product").Preload("UOM").Find(&adjustmentDetails, " adjustment_id = ? ", adjustmentID)

	return adjustmentDetails
}

// GetAdjustmentDetailPage ...
func GetAdjustmentDetailPage(param dto.FilterAdjustmentDetail, offset, limit int) ([]dbmodels.AdjustmentDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var adjustmentDetails []dbmodels.AdjustmentDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&adjustmentDetails).Error
		if err != nil {
			return adjustmentDetails, 0, err
		}
		return adjustmentDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysAdjustmentDetails(db, offset, limit, &adjustmentDetails, param.AdjustmentID, errQuery)
	go AsyncQueryCountsAdjustmentDetails(db, &total, param.AdjustmentID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return adjustmentDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return adjustmentDetails, 0, resErrCount
	}
	return adjustmentDetails, total, nil
}

// AsyncQueryCountsAdjustmentDetails ...
func AsyncQueryCountsAdjustmentDetails(db *gorm.DB, total *int, adjustmentID int64, offset int, limit int, resChan chan error) {

	var err error

	err = db.Model(&dbmodels.AdjustmentDetail{}).Where("adjustment_id = ?", adjustmentID).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysAdjustmentDetails ...
func AsyncQuerysAdjustmentDetails(db *gorm.DB, offset int, limit int, adjustmentDetails *[]dbmodels.AdjustmentDetail, adjustmentID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).Preload("Product").Preload("UOM").Order("id desc").Find(&adjustmentDetails, "adjustment_id = ? ", adjustmentID).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("adj--> ", adjustmentDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//SaveAdjustmentDetail ...
func SaveAdjustmentDetail(adjustmentDetail *dbmodels.AdjustmentDetail) (errCode string, errDesc string) {

	fmt.Println(" Update Adjustment Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Save(&adjustmentDetail); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", adjustmentDetail.ID)
	return

}

// DeleteAdjustmentDetailById ...
func DeleteAdjustmentDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Adjustment Detail  ------------------------------------------  %v ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Where("id = ? ", id).Delete(dbmodels.AdjustmentDetail{}); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", id)
	return

}
