package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

//SaveStockMutation ...
func SaveStockMutation(stockMutation *dbmodels.StockMutation) (errCode string, errDesc string, id int64, status int8) {

	fmt.Println(" Update  StockMutation numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Save(&stockMutation)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, stockMutation.ID, stockMutation.Status
}

// SaveStockMutationApprove ...
func SaveStockMutationApprove(stockMutation *dbmodels.StockMutation) (errCode string, errDesc string) {

	fmt.Println(" Approve StockMutation ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var total float32
	total = 0
	stockMutationDetails := GetAllDataStockMutationDetail(stockMutation.ID)
	for idx, stockMutationDetail := range stockMutationDetails {
		fmt.Println("idx -> ", idx)

		productHpp := getHppByProductID(stockMutationDetail.ProductID)
		productId := stockMutationDetail.ProductID

		curProduct, _, _ := FindProductByID(productId)

		// kurangkan stock source gudang
		fmt.Println("GET stock warehouse source")
		curStockQtySourceWh, _, _ := GetStockByProductAndWarehouse(productId, stockMutation.WarehouseSourceID)
		fmt.Println("UPDATE stock warehouse source ON HAND = ", curStockQtySourceWh.Qty)
		UpdateStockProductByID(productId, curStockQtySourceWh.Qty-stockMutationDetail.Qty, stockMutation.WarehouseSourceID)

		// tambahkan stock dest
		fmt.Println("GET stock warehouse dest")
		curStockQtyDestWh, _, _ := GetStockByProductAndWarehouse(productId, stockMutation.WarehouseDestID)
		fmt.Println("Update stock warehouse dest = ", curStockQtyDestWh.Qty)
		UpdateStockProductByID(productId, curStockQtyDestWh.Qty+stockMutationDetail.Qty, stockMutation.WarehouseDestID)

		fmt.Println("start history stock")
		var historySource dbmodels.HistoryStock
		historySource.Code = curProduct.Code
		historySource.Name = curProduct.Name
		historySource.Debet = 0
		historySource.Kredit = stockMutationDetail.Qty
		historySource.Saldo = curStockQtySourceWh.Qty - stockMutationDetail.Qty
		historySource.Description = "Mutation "
		historySource.TransDate = stockMutation.StockMutationDate
		historySource.ReffNo = stockMutation.StockMutationNo
		historySource.Price = 0
		historySource.Hpp = productHpp
		historySource.LastUpdateBy = stockMutation.LastUpdateBy
		historySource.LastUpdate = util.GetCurrDate()
		historySource.Disc1 = 0
		historySource.Total = 0
		historySource.WarehouseID = stockMutation.WarehouseSourceID
		db.Save(&historySource)

		var historyDest dbmodels.HistoryStock
		historyDest.Code = curProduct.Code
		historyDest.Name = curProduct.Name
		historyDest.Debet = stockMutationDetail.Qty
		historyDest.Kredit = 0
		historyDest.Saldo = curStockQtyDestWh.Qty + stockMutationDetail.Qty
		historyDest.Description = "Mutation "
		historyDest.TransDate = stockMutation.StockMutationDate
		historyDest.ReffNo = stockMutation.StockMutationNo
		historyDest.Price = 0
		historyDest.Hpp = productHpp
		historyDest.LastUpdateBy = stockMutation.LastUpdateBy
		historyDest.LastUpdate = util.GetCurrDate()
		historyDest.Disc1 = 0
		historyDest.Total = 0
		historyDest.WarehouseID = stockMutation.WarehouseDestID
		db.Save(&historyDest)

		// Update hpp

		fmt.Println("update hpp -> ", productHpp)
		stockMutationDetail.Hpp = productHpp
		stockMutationDetail.LastUpdate = util.GetCurrDate()
		stockMutationDetail.LastUpdateBy = dto.CurrUser
		db.Save(&stockMutationDetail)

		total = total + (productHpp)
	}

	db.Debug().LogMode(true)

	stockMutation.Total = total
	stockMutation.LastUpdateBy = dto.CurrUser
	stockMutation.LastUpdate = util.GetCurrDate()
	stockMutation.Status = 20
	r := db.Save(&stockMutation)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update stock mutation ", errDesc)
		return
	}

	tx.Commit()
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetStockMutationByStockMutationNo ...
func GetStockMutationByStockMutationNo(stockMutationNo string) (dbmodels.StockMutation, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	order := dbmodels.StockMutation{}

	err := db.Preload("WarehouseSource").Preload("WarehouseDest").Where(" mutation_no = ?  ", stockMutationNo).First(&order).Error

	return order, err

}

// GetStockMutationByStockMutationId ...
func GetStockMutationById(stockMutationID int64) (dbmodels.StockMutation, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	stockMutation := dbmodels.StockMutation{}

	err := db.Preload("WarehouseSource").Preload("WarehouseDest").Where(" id = ?  ", stockMutationID).First(&stockMutation).Error
	return stockMutation, err

}

// GetStockMutationPage ...
func GetStockMutationPage(param dto.FilterStockMutation, offset, limit, internalStatus int) ([]dbmodels.StockMutation, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stockMutations []dbmodels.StockMutation
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&stockMutations).Error
		if err != nil {
			return stockMutations, 0, err
		}
		return stockMutations, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysStockMutations(db, offset, limit, internalStatus, &stockMutations, param, errQuery)
	go AsyncQueryCountsStockMutations(db, &total, internalStatus, &stockMutations, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return stockMutations, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return stockMutations, 0, resErrCount
	}
	return stockMutations, total, nil
}

func getParamMutation(param dto.FilterStockMutation, status int) (stockMutationNumber string, byStatus bool) {

	stockMutationNumber = param.StockMutationNumber
	if stockMutationNumber == "" {
		stockMutationNumber = "%"
	} else {
		stockMutationNumber = "%" + param.StockMutationNumber + "%"
	}

	byStatus = true
	if status == -1 {
		byStatus = false
	}

	return
}

// AsyncQueryCountsStockMutations ...
func AsyncQueryCountsStockMutations(db *gorm.DB, total *int, status int, stockMutation *[]dbmodels.StockMutation, param dto.FilterStockMutation, resChan chan error) {

	stockMutationNumber, byStatus := getParamMutation(param, status)
	fmt.Println("ISI  mutation Number ", stockMutationNumber, "  status ", status, " fill status ", byStatus)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&stockMutation).Where(" ( (status = ?) or ( not ?) ) AND COALESCE(mutation_no, '') ilike ? AND mutation_date between ? and ?  ", status, byStatus, stockMutationNumber, param.StartDate, param.EndDate).Count(&*total).Error
	} else {
		err = db.Model(&stockMutation).Where(" ( (status = ?) or ( not ?) ) AND COALESCE(mutation_no,'') ilike ? ", status, byStatus, stockMutationNumber).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysStockMutations ...
func AsyncQuerysStockMutations(db *gorm.DB, offset int, limit int, status int, stockMutations *[]dbmodels.StockMutation, param dto.FilterStockMutation, resChan chan error) {

	var err error
	stockMutationNumber, byStatus := getParamMutation(param, status)
	fmt.Println("ISI  mutation Number ", stockMutationNumber, "  status ", status, " fill status ", byStatus)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("WarehouseSource").Preload("WarehouseDest").Order("mutation_date DESC").Offset(offset).Limit(limit).Find(&stockMutations, " ( ( status = ?) or ( not ?) ) AND COALESCE(mutation_no, '') ilike ? AND mutation_date between ? and ?   ", status, byStatus, stockMutationNumber, param.StartDate, param.EndDate).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("WarehouseSource").Preload("WarehouseDest").Order("mutation_date DESC").Find(&stockMutations, " ( ( status = ?) or ( not ?) ) AND COALESCE(mutation_no,'') ilike ?  ", status, byStatus, stockMutationNumber).Error
		if err != nil {
			fmt.Println("error --> ", err)
		}
		fmt.Println("order--> ", stockMutations)

	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//RejectStockMutation ...
func RejectStockMutation(stockMutation *dbmodels.StockMutation) (errCode string, errDesc string) {

	fmt.Println(" Reject StockMutation ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.StockMutation{}).Where("id =?", stockMutation.ID).Update(dbmodels.StockMutation{Status: 30})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
