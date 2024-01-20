package database

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

//SaveHistory ...
func SaveHistory(history dbmodels.HistoryStock) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Save(&history); r.Error != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error save data to DB"
		return res
	}

	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG
	return res
}

// GetHistoryPage ...
func GetHistoryPage(param dto.FilterHistoryStock, offset, limit int) ([]dbmodels.HistoryStock, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var historyStocks []dbmodels.HistoryStock
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&historyStocks).Error
		if err != nil {
			return historyStocks, 0, err
		}
		return historyStocks, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysHistoryStock(db, offset, limit, &historyStocks, param, errQuery)
	go AsyncQueryCountsHistoryStocks(db, &total, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return historyStocks, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return historyStocks, 0, resErrCount
	}
	fmt.Println("Total Rec search = ", total)
	return historyStocks, total, nil
}

func AsyncQueryCountsHistoryStocks(db *gorm.DB, total *int, param dto.FilterHistoryStock, resChan chan error) {
	// func AsyncQueryCountsHistoryStocks(db *gorm.DB, total *int, historyStocks *[]dbmodels.HistoryStock, param dto.FilterHistoryStock, resChan chan error) {

	var err error
	// err = db.Model(&historyStocks).Where(" DATE(trans_date)  between DATE(?)  and DATE(?)  and warehouse_id = ? and code = ?  ", param.StartDate, param.EndDate, param.WarehouseID, param.ProductCode).Count(&*total).Error
	err = db.Model(dbmodels.HistoryStock{}).Where(" DATE(trans_date)  between DATE(?)  and DATE(?)  and warehouse_id = ? and code = ?  ", param.StartDate, param.EndDate, param.WarehouseID, param.ProductCode).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysHistoryStock ...
func AsyncQuerysHistoryStock(db *gorm.DB, offset int, limit int, historyStocks *[]dbmodels.HistoryStock, param dto.FilterHistoryStock, resChan chan error) {

	var err error

	fmt.Println(" Date [", param.StartDate, "]  sd [", param.EndDate, "] product [", param.ProductCode, "] warehouse  [", param.WarehouseID, "] ")

	err = db.Order("id DESC").Offset(offset).Limit(limit).Find(&historyStocks, " DATE(trans_date) between  DATE(?)  and DATE(?)  and warehouse_id = ? and code = ?  ", param.StartDate, param.EndDate, param.WarehouseID, param.ProductCode).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}
