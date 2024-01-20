package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// GetStockByProductAndWarehouse ...
func GetStockByProductAndWarehouse(productID, warehouseID int64) (dbmodels.Stock, string, string) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stock dbmodels.Stock
	err := db.Where("warehouse_id = ? and product_id = ? ", warehouseID, productID).First(&stock).Error

	if err != nil {
		return stock, constants.ERR_CODE_81, constants.ERR_CODE_81_MSG + " " + err.Error()
	}
	return stock, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// GetStockByProduct ...
func GetStockByProductPage(productID int64, offset int, limit int) ([]dbmodels.Stock, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stocks []dbmodels.Stock
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&stocks).Error
		if err != nil {
			return stocks, 0, err
		}
		return stocks, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQueryStock(db, offset, limit, &stocks, productID, errQuery)
	go AsyncQueryCountStock(db, &total, &stocks, productID, errCount)

	resErrCount := <-errCount
	resErrQuery := <-errQuery

	wg.Done()

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return stocks, 0, resErrCount
	}

	if resErrQuery != nil {
		return stocks, 0, resErrQuery
	}

	return stocks, total, nil
}

// AsyncQueryStock ...
func AsyncQueryStock(db *gorm.DB, offset int, limit int, stocks *[]dbmodels.Stock, productID int64, resChan chan error) {

	var err error
	err = db.Set("gorm:auto_preload", true).Order("warehouse_id ASC").Offset(offset).Limit(limit).Find(&stocks, " product_id = ?   ", productID).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQueryCountStock ...
func AsyncQueryCountStock(db *gorm.DB, total *int, stock *[]dbmodels.Stock, productID int64, resChan chan error) {

	err := db.Model(&stock).Where(" product_id = ?  ", productID).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// SaveStock ...
func SaveStock(productID, warehouseID int64) (dbmodels.Stock, string, string) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stock dbmodels.Stock
	stock.ProductID = productID
	stock.LastUpdateBy = dto.CurrUser
	stock.LastUpdate = time.Now()
	stock.Qty = 0
	stock.WarehouseID = warehouseID

	err := db.Save(&stock).Error

	if err != nil {
		return stock, constants.ERR_CODE_81, constants.ERR_CODE_81_MSG + " " + err.Error()
	}
	return stock, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}
