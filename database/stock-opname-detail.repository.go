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

// GetStockStockOpnameDetailPage ...
func GetStockOpnameDetailPage(param dto.FilterStockOpname, offset, limit int) ([]dbmodels.StockOpnameDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stockOpnameDetails []dbmodels.StockOpnameDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&stockOpnameDetails).Error
		if err != nil {
			return stockOpnameDetails, 0, err
		}
		return stockOpnameDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysStockOpnameDetails(db, offset, limit, &stockOpnameDetails, param.StockOpnameID, errQuery)
	go AsyncQueryCountsStockOpnameDetails(db, &total, param.StockOpnameID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return stockOpnameDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return stockOpnameDetails, 0, resErrCount
	}
	return stockOpnameDetails, total, nil
}

// AsyncQueryCountsStockOpnameDetails ...
func AsyncQueryCountsStockOpnameDetails(db *gorm.DB, total *int, stockOpnameID int64, offset int, limit int, resChan chan error) {

	var err error
	err = db.Model(&dbmodels.StockOpnameDetail{}).Where("stock_opname_id = ?", stockOpnameID).Count(total).Error
	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysStockOpnameDetails ...
func AsyncQuerysStockOpnameDetails(db *gorm.DB, offset int, limit int, stockOpnameDetails *[]dbmodels.StockOpnameDetail, stockOpnameID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).Preload("Product").Preload("UOM").Order("id asc").Find(&stockOpnameDetails, " stock_opname_id = ? ", stockOpnameID).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("order--> ", stockOpnameDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// GetAllDataStockOpnameDetail ...
func GetAllDataStockOpnameDetail(stockOpnameID int64) []dbmodels.StockOpnameDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var stockOpnameDetails []dbmodels.StockOpnameDetail

	db.Preload("Product").Preload("UOM").Find(&stockOpnameDetails, " stock_opname_id = ?  ", stockOpnameID)

	return stockOpnameDetails
}

//SaveStockOpnameDetail ...
func SaveStockOpnameDetail(stockOpnameDetail *dbmodels.StockOpnameDetail) (errCode string, errDesc string) {

	fmt.Println(" Update Stock Mutation  Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", stockOpnameDetail.ID)

	stockOpnameDetail.Hpp = getHppByProductID(stockOpnameDetail.ProductID)
	stockOpnameDetail.LastUpdateBy = dto.CurrUser
	stockOpnameDetail.LastUpdate = util.GetCurrDate()
	if r := db.Save(&stockOpnameDetail); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}

func DeleteStockOpnameDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Stock Mutation Detail  ------------------------------------------  ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", id)

	if r := db.Where("id = ? ", id).Delete(dbmodels.StockOpnameDetail{}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}

//UpdateQtyStockOpnameDetail ...
func UpdateQtyStockOpnameDetail(stockOpnameDetailId int64, qty int64) (errCode string, errDesc string) {

	fmt.Println(" Update Qty stock Mutation Detail  -- ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("id = %v, qty = %v", stockOpnameDetailId, qty)

	if r := db.Model(&dbmodels.StockOpnameDetail{}).
		Where("id = ?", stockOpnameDetailId).
		Update(dbmodels.StockOpnameDetail{
			Qty:          qty,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}
