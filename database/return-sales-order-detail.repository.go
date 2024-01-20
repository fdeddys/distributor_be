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

// GetAllSalesOrderReturnDetail ...
func GetAllSalesOrderReturnDetail(orderReturnID int64) []dbmodels.ReturnSalesOrderDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var orderReturnDetails []dbmodels.ReturnSalesOrderDetail

	db.Preload("Product").Preload("UOM").Find(&orderReturnDetails, " return_sales_order_id = ?  ", orderReturnID)

	return orderReturnDetails
}

// GetReturnOrderDetailPage ...
func GetReturnOrderDetailPage(param dto.FilterOrderReturnDetail, offset, limit int) ([]dbmodels.ReturnSalesOrderDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var orderReturnDetails []dbmodels.ReturnSalesOrderDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&orderReturnDetails).Error
		if err != nil {
			return orderReturnDetails, 0, err
		}
		return orderReturnDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysOrderReturnDetails(db, offset, limit, &orderReturnDetails, param.OrderReturnID, errQuery)
	go AsyncQueryCountsOrderReturnDetails(db, &total, param.OrderReturnID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return orderReturnDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return orderReturnDetails, 0, resErrCount
	}
	return orderReturnDetails, total, nil
}

// AsyncQueryCountsOrderDetails ...
func AsyncQueryCountsOrderReturnDetails(db *gorm.DB, total *int, returnOrderID int64, offset int, limit int, resChan chan error) {

	var err error

	err = db.Model(&dbmodels.ReturnSalesOrderDetail{}).Offset(offset).Where("return_sales_order_id = ?", returnOrderID).Count(total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrderReturnDetails ...
func AsyncQuerysOrderReturnDetails(db *gorm.DB, offset int, limit int, returnOrderDetails *[]dbmodels.ReturnSalesOrderDetail, returnOrderID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).Preload("Product").Preload("UOM").Order("id asc").Find(&returnOrderDetails, " return_sales_order_id = ? ", returnOrderID).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("order--> ", returnOrderDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//SaveReturnOrderDetail ...
func SaveReturnOrderDetail(returnOrderDetail *dbmodels.ReturnSalesOrderDetail) (errCode string, errDesc string) {

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

// DeleteReturnOrderDetailById ...
func DeleteReturnOrderDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Return Detail  ------------------------------------------ ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Where("id = ? ", id).Delete(dbmodels.ReceiveDetail{}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = r.Error.Error()
	}

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", id)
	return

}

//UpdateQtyReturnSalesOrderDetail ...
func UpdateQtyReturnSalesOrderDetail(returnSaleOrderDetailId int64, qty int64) (errCode string, errDesc string) {

	fmt.Println(" Update Qty return Sales Order Detail  -- ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("id = %v, qty = %v", returnSaleOrderDetailId, qty)

	if r := db.Model(&dbmodels.ReturnSalesOrderDetail{}).
		Where("id = ?", returnSaleOrderDetailId).
		Update(dbmodels.ReturnSalesOrderDetail{
			Qty:          qty,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		}); r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return

}
