package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

//SaveAdjustment ...
func SaveAdjustment(adjustment *dbmodels.Adjustment) (errCode string, errDesc string, id int64, status int8) {

	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Save(&adjustment)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, adjustment.ID, adjustment.Status
}

// SaveAdjustmentApprove ...
func SaveAdjustmentApprove(adjustment *dbmodels.Adjustment) (errCode string, errDesc string) {

	fmt.Println(" Approve Receiving ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// update stock
	// update history stock
	// hitung ulang
	var total float32
	total = 0
	adjustmentDetails := GetAllDataDetailAdjustment(adjustment.ID)
	for idx, adjustmentDetail := range adjustmentDetails {
		fmt.Println("idx -> ", idx)

		product, errCodeProd, errDescProd := FindProductByID(adjustmentDetail.ProductID)
		if errCodeProd != constants.ERR_CODE_00 {
			tx.Rollback()
			return errCodeProd, errDescProd
		}

		checkStock, _, _ := GetStockByProductAndWarehouse(product.ID, adjustment.WarehouseID)

		curQty := checkStock.Qty
		curHpp := product.Hpp
		updateQty := int64(0)

		var historyStock dbmodels.HistoryStock
		if adjustmentDetail.Qty > 0 {
			historyStock.Debet = adjustmentDetail.Qty
			historyStock.Kredit = 0
			// updateQty = curQty + adjustmentDetail.Qty
		} else {
			historyStock.Debet = 0
			historyStock.Kredit = -1 * adjustmentDetail.Qty
		}
		updateQty = curQty + adjustmentDetail.Qty
		historyStock.Saldo = updateQty
		historyStock.Code = product.Code
		historyStock.Description = "Adjustment"
		historyStock.Hpp = curHpp
		historyStock.Name = product.Name
		historyStock.Price = curHpp
		historyStock.ReffNo = adjustment.AdjustmentNo
		historyStock.TransDate = adjustment.AdjustmentDate
		historyStock.WarehouseID = adjustment.WarehouseID
		historyStock.LastUpdate = time.Now()
		historyStock.LastUpdateBy = dto.CurrUser
		historyStock.Total = float32(adjustmentDetail.Qty * int64(adjustmentDetail.Hpp))

		UpdateStockAndHppProductByID(adjustmentDetail.ProductID, adjustment.WarehouseID, updateQty, curHpp)
		SaveHistory(historyStock)
		total = total + (curHpp * float32(adjustmentDetail.Qty))
	}

	db.Debug().LogMode(true)
	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	adjustment.Total = total
	adjustment.LastUpdateBy = dto.CurrUser
	adjustment.LastUpdate = time.Now()
	adjustment.Status = 20
	r := db.Save(&adjustment)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	// fmt.Println("Order [database]=> order id", order.OrderNo)

	tx.Commit()
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetAdjustmentPage ...
func GetAdjustmentPage(param dto.FilterAdjustment, offset, limit, internalStatus int) ([]dbmodels.Adjustment, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var adjustments []dbmodels.Adjustment
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&adjustments).Error
		if err != nil {
			return adjustments, 0, err
		}
		return adjustments, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysAdjustments(db, offset, limit, internalStatus, &adjustments, param, errQuery)
	go AsyncQueryCountsAdjustments(db, &total, internalStatus, &adjustments, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return adjustments, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return adjustments, 0, resErrCount
	}
	return adjustments, total, nil
}

// AsyncQueryCountsAdjustments ...
func AsyncQueryCountsAdjustments(db *gorm.DB, total *int, status int, orders *[]dbmodels.Adjustment, param dto.FilterAdjustment, resChan chan error) {

	adjustmentNumber, byStatus := getParamAdjustment(param, status)

	fmt.Println(" Rec Number ", adjustmentNumber, "  status ", status, " fill status ", byStatus)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&orders).Where(" ( (status = ?) or ( not ?) ) AND  COALESCE(adjustment_no, '') ilike ? AND adjustment_date between ? and ?  ", status, byStatus, adjustmentNumber, param.StartDate, param.EndDate).Count(&*total).Error
	} else {
		err = db.Model(&orders).Where(" ( (status = ?) or ( not ?) ) AND COALESCE(adjustment_no,'') ilike ? ", status, byStatus, adjustmentNumber).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysAdjustments ...
func AsyncQuerysAdjustments(db *gorm.DB, offset int, limit int, status int, adjustments *[]dbmodels.Adjustment, param dto.FilterAdjustment, resChan chan error) {

	var err error

	adjustmentNumber, byStatus := getParamAdjustment(param, status)

	fmt.Println(" Adjustment no ", adjustmentNumber, "  status ", status, " fill status ", byStatus)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Order("adjustment_date DESC, id desc").Offset(offset).Limit(limit).Find(&adjustments, " ( ( status = ?) or ( not ?) ) AND COALESCE(adjustment_no, '') ilike ? AND adjustment_date between ? and ?   ", status, byStatus, adjustmentNumber, param.StartDate, param.EndDate).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Order("adjustment_date DESC, id desc").Offset(offset).Limit(limit).Find(&adjustments, " ( ( status = ?) or ( not ?) ) AND COALESCE(adjustment_no,'') ilike ?  ", status, byStatus, adjustmentNumber).Error
		if err != nil {
			fmt.Println("adjustment --> ", err)
		}
		fmt.Println("adjustment--> ", adjustments)

	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

func getParamAdjustment(param dto.FilterAdjustment, status int) (adjustmentNumber string, byStatus bool) {

	adjustmentNumber = param.AdjustmentNumber
	if adjustmentNumber == "" {
		adjustmentNumber = "%"
	} else {
		adjustmentNumber = "%" + param.AdjustmentNumber + "%"
	}

	byStatus = true
	if status == -1 {
		byStatus = false
	}

	return
}

// GetAdjustmentByAdjustmentID ...
func GetAdjustmentByAdjustmentID(adjustmentID int64) (dbmodels.Adjustment, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	adjustment := dbmodels.Adjustment{}

	err := db.Where(" id = ?  ", adjustmentID).First(&adjustment).Error

	return adjustment, err

}

//RejectAdjustment ...
func RejectAdjustment(adjustment *dbmodels.Adjustment) (errCode string, errDesc string) {

	fmt.Println(" Reject Adjustment numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.Adjustment{}).Where("id =?", adjustment.ID).Update(dbmodels.Adjustment{Status: 20})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
