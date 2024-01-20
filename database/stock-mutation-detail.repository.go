package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetStockMutationDetailPage ...
func GetStockMutationDetailPage(param dto.FilterStockMutation, offset, limit int) ([]dbmodels.StockMutationDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stockMutationDetails []dbmodels.StockMutationDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&stockMutationDetails).Error
		if err != nil {
			return stockMutationDetails, 0, err
		}
		return stockMutationDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysStockMutationDetails(db, offset, limit, &stockMutationDetails, param.StockMutationID, errQuery)
	go AsyncQueryCountsStockMutationDetails(db, &total, param.StockMutationID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return stockMutationDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return stockMutationDetails, 0, resErrCount
	}
	return stockMutationDetails, total, nil
}

// AsyncQueryCountsStockMutationDetails ...
func AsyncQueryCountsStockMutationDetails(db *gorm.DB, total *int, stockMutationID int64, offset int, limit int, resChan chan error) {

	var err error
	err = db.Model(&dbmodels.StockMutationDetail{}).Offset(offset).Where("mutation_id = ?", stockMutationID).Count(total).Error
	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysStockMutationDetails ...
func AsyncQuerysStockMutationDetails(db *gorm.DB, offset int, limit int, stockMutationDetails *[]dbmodels.StockMutationDetail, stockMutationID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).Preload("Product").Preload("UOM").Order("id asc").Find(&stockMutationDetails, " mutation_id = ? ", stockMutationID).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("order--> ", stockMutationDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// GetAllDataStockMutationDetail ...
func GetAllDataStockMutationDetail(stockMutationID int64) []dbmodels.StockMutationDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var stockMutationDetails []dbmodels.StockMutationDetail

	db.Preload("Product").Preload("UOM").Find(&stockMutationDetails, " mutation_id = ?  ", stockMutationID)

	return stockMutationDetails
}

//SaveStockMutationDetail ...
func SaveStockMutationDetail(stockMutationDetail *dbmodels.StockMutationDetail) (errCode string, errDesc string) {

	fmt.Println(" Update Stock Mutation  Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", stockMutationDetail.ID)

	stockMutationDetail.Hpp = getHppByProductID(stockMutationDetail.ProductID)
	stockMutationDetail.LastUpdateBy = dto.CurrUser
	stockMutationDetail.LastUpdate = util.GetCurrDate()
	if r := db.Save(&stockMutationDetail); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}

func getHppByProductID(productID int64) float32 {

	db := GetDbCon()
	db.Debug().LogMode(true)

	product, errCode, _ := FindProductByID(productID)

	if errCode == constants.ERR_CODE_00 {
		return product.Hpp
	}
	return 0
}

func DeleteStockMutationDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Stock Mutation Detail  ------------------------------------------  ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", id)

	if r := db.Where("id = ? ", id).Delete(dbmodels.StockMutationDetail{}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}

//UpdateQtyStockMutationDetail ...
func UpdateQtyStockMutationDetail(stockMutationDetailId int64, qty int64) (errCode string, errDesc string) {

	fmt.Println(" Update Qty stock Mutation Detail  -- ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("id = %v, qty = %v", stockMutationDetailId, qty)

	if r := db.Model(&dbmodels.StockMutationDetail{}).
		Where("id = ?", stockMutationDetailId).
		Update(dbmodels.StockMutationDetail{
			Qty:          qty,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}
