package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

//SaveReceive ...
func SaveReceive(receive *dbmodels.Receive) (errCode string, errDesc string, id int64, status int8) {

	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Save(&receive)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, receive.ID, receive.Status
}

// SaveReceiveApprove ...
func SaveReceiveApprove(receive *dbmodels.Receive) (errCode string, errDesc string) {

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
	var subtotal float32
	var grandTotal float32
	total = 0
	grandTotal = 0
	subtotal = 0
	receiveDetails := GetAllDataDetailReceive(receive.ID)
	for idx, receiveDetail := range receiveDetails {
		fmt.Println("idx -> ", idx)
		addNewStock := false
		product, errCodeProd, errDescProd := FindProductByID(receiveDetail.ProductID)
		if errCodeProd != constants.ERR_CODE_00 {
			tx.Rollback()
			return errCodeProd, errDescProd
		}

		checkStock, errCode, _ := GetStockByProductAndWarehouse(product.ID, receive.WarehouseID)

		if errCode != constants.ERR_CODE_00 {
			addNewStock = true
		}

		curQty := int64(0)
		updateQty := int64(0)
		newHpp := float32(0)
		if addNewStock {
			updateQty = receiveDetail.Qty
			newHpp = receiveDetail.Price
		} else {
			curQty = checkStock.Qty
			updateQty = curQty + receiveDetail.Qty
			newHpp = reCalculateHpp(product.Hpp, checkStock.Qty, receiveDetail.Price, receiveDetail.Qty, receiveDetail.Disc1, receiveDetail.Disc2)
		}
		// curQty := checkStock.Qty

		var historyStock dbmodels.HistoryStock
		historyStock.Code = product.Code
		historyStock.Description = "Receive"
		historyStock.Hpp = newHpp
		historyStock.Name = product.Name
		historyStock.Price = receiveDetail.Price
		historyStock.ReffNo = receive.ReceiveNo
		historyStock.TransDate = receive.ReceiveDate
		historyStock.Debet = receiveDetail.Qty
		historyStock.Kredit = 0
		historyStock.Saldo = updateQty
		historyStock.LastUpdate = time.Now()
		historyStock.LastUpdateBy = dto.CurrUser
		historyStock.Disc1 = receiveDetail.Disc1
		historyStock.Disc2 = receiveDetail.Disc2
		historyStock.WarehouseID = receive.WarehouseID

		//total -= (receiveDetail.Price * float32(receiveDetail.Qty) * ((100 - receiveDetail.Disc1) / 100))
		total = receiveDetail.Price * float32(receiveDetail.Qty)
		total = total * ((100 - receiveDetail.Disc1) / 100)
		total = total * ((100 - receiveDetail.Disc2) / 100)
		historyStock.Total = total
		fmt.Println("total -> ", total)

		if addNewStock {
			AddnewStockAndHppProductByID(receiveDetail.ProductID, receive.WarehouseID, updateQty, newHpp)
		} else {
			UpdateStockAndHppProductByID(receiveDetail.ProductID, receive.WarehouseID, updateQty, newHpp)
		}
		db.Save(&historyStock)
		subtotal += total
	}

	db.Debug().LogMode(true)
	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	grandTotal = subtotal
	curTax := float32(0)
	if receive.Tax != 0 {
		curTax = float32(math.Round(float64(grandTotal) * float64(receive.Tax) / 100))
	}
	grandTotal += curTax
	receive.GrandTotal = grandTotal
	receive.Total = subtotal
	receive.LastUpdateBy = dto.CurrUser
	receive.LastUpdate = time.Now()
	receive.Status = 20
	r := db.Save(&receive)
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

func reCalculateHpp(hpp1 float32, qty1 int64, price2 float32, qty2 int64, disc1, disc2 float32) float32 {

	newHpp := price2
	newHpp = newHpp * ((100 - disc1) / 100)
	newHpp = newHpp * ((100 - disc2) / 100)

	totalRp := (hpp1 * float32(qty1)) + (price2 * float32(qty2))
	totalQty := qty1 + qty2
	return (totalRp / float32(totalQty))
}

// GetReceivePage ...
func GetReceivePage(param dto.FilterReceive, offset, limit, internalStatus int) ([]dbmodels.Receive, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var receives []dbmodels.Receive
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&receives).Error
		if err != nil {
			return receives, 0, err
		}
		return receives, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysReceives(db, offset, limit, internalStatus, &receives, param, errQuery)
	go AsyncQueryCountsReceives(db, &total, internalStatus, &receives, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return receives, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return receives, 0, resErrCount
	}
	return receives, total, nil
}

// AsyncQueryCountsReceives ...
func AsyncQueryCountsReceives(db *gorm.DB, total *int, status int, orders *[]dbmodels.Receive, param dto.FilterReceive, resChan chan error) {

	receiveNumber, byStatus, suppName, pono := getParamReceive(param, status)

	fmt.Println(" Rec Number ", receiveNumber, "  status ", status, " fill status ", byStatus)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.
			Model(&orders).
			Joins("inner join supplier on supplier.id = receive.supplier_id and supplier.name ilike ? ", suppName).
			Preload("Supplier", " name ilike ? ", suppName).
			Where(" ( (receive.status = ?) or ( not ?) ) AND  COALESCE(receive_no, '') ilike ? AND receive_date between ? and ?  AND COALESCE(po_no, '') ilike ?  ", status, byStatus, receiveNumber, param.StartDate, param.EndDate, pono).Count(&*total).
			Error
	} else {
		err = db.
			Model(&orders).
			Joins("inner join supplier on supplier.id = receive.supplier_id and supplier.name ilike ? ", suppName).
			Preload("Supplier", " name ilike ? ", suppName).Where(" ( (receive.status = ?) or ( not ?) ) AND COALESCE(receive_no,'') ilike ? AND COALESCE(po_no, '') ilike ?  ", status, byStatus, receiveNumber, pono).Count(&*total).
			Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysReceives ...
func AsyncQuerysReceives(db *gorm.DB, offset int, limit int, status int, receives *[]dbmodels.Receive, param dto.FilterReceive, resChan chan error) {

	var err error

	receiveNumber, byStatus, suppName, pono := getParamReceive(param, status)

	fmt.Println(" Receive no ", receiveNumber, "  status ", status, " fill status ", byStatus)

	fmt.Println("isi dari filter => receive [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter =>masukXX [", param.StartDate, '-', param.EndDate, "] ")
		// err = db.Order("receive_date DESC, id DESC").Offset(offset).Limit(limit).Preload("Supplier", " name like ? ", suppName).Find(&receives, " ( ( status = ?) or ( not ?) ) AND COALESCE(receive_no, '') ilike ? AND receive_date between ? and ?   ", status, byStatus, receiveNumber, param.StartDate, param.EndDate).Error
		// err = db.Order("receive_date DESC, id DESC").Offset(offset).Limit(limit).Preload("Supplier", func(db *gorm.DB) *gorm.DB {
		// 	return db.Where(" name ilike ? ", suppName)
		// }).Find(&receives, " ( ( status = ?) or ( not ?) ) AND COALESCE(receive_no, '') ilike ? AND receive_date between ? and ?   ", status, byStatus, receiveNumber, param.StartDate, param.EndDate).Error
		err = db.
			Joins("inner join supplier on supplier.id = receive.supplier_id and supplier.name ilike ? ", suppName).
			Order("receive_date DESC, id DESC").
			Preload("Supplier").
			Offset(offset).
			Limit(limit).
			Find(&receives, " ( ( receive.status = ?) or ( not ?) ) AND COALESCE(receive_no, '') ilike ? AND receive_date between ? and ?  AND COALESCE(po_no, '') ilike ? ", status, byStatus, receiveNumber, param.StartDate, param.EndDate, pono).
			// Find(&receives).
			Error

	} else {
		fmt.Println("isi dari kosong ")
		err = db.
			Joins("inner join supplier on supplier.id = receive.supplier_id and supplier.name ilike ? ", suppName).
			Order("receive_date DESC, id DESC").
			Offset(offset).Limit(limit).
			Preload("Supplier", " name ilike ? ", suppName).
			Find(&receives, " ( ( receive.status = ?) or ( not ?) ) AND COALESCE(receive_no,'') ilike ? AND COALESCE(po_no, '') ilike ?  ", status, byStatus, receiveNumber, pono).
			Error
		if err != nil {
			fmt.Println("receive --> ", err)
		}
		fmt.Println("receive--> ", receives)

	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

func getParamReceive(param dto.FilterReceive, status int) (receiveNumber string, byStatus bool, supplierName string, pono string) {

	receiveNumber = param.ReceiveNumber
	if receiveNumber == "" {
		receiveNumber = "%"
	} else {
		receiveNumber = "%" + param.ReceiveNumber + "%"
	}

	byStatus = true
	if status == -1 {
		byStatus = false
	}

	supplierName = param.SupplierName
	if supplierName == "" {
		supplierName = "%"
	} else {
		supplierName = "%" + param.SupplierName + "%"
	}

	pono = param.PurchaseOrderNo
	if pono == "" {
		pono = "%"
	} else {
		pono = "%" + param.PurchaseOrderNo + "%"
	}

	return
}

// GetReceiveByReceiveID ...
func GetReceiveByReceiveID(receiveID int64) (dbmodels.Receive, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	receive := dbmodels.Receive{}

	err := db.Preload("Supplier").Preload("Warehouse").Where(" id = ?  ", receiveID).First(&receive).Error

	return receive, err

}

//RejectReceive ...
func RejectReceive(receive *dbmodels.Receive) (errCode string, errDesc string) {

	fmt.Println(" Reject Receive numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.Receive{}).Where("id =?", receive.ID).Update(dbmodels.Receive{Status: 30})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// RemovePO ...
func RemovePO(receive *dbmodels.Receive, removeItem bool) (errCode string, errDesc string) {

	fmt.Println("Remove PO Receiving ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// remove PO numb from Receive
	// update PO status

	currPO := dbmodels.PurchaseOrder{}
	tx.First(&currPO, "po_no = ? ", receive.PoNo)
	currPO.Status = 20
	tx.Save(currPO)

	currReceive := dbmodels.Receive{}
	tx.Find(&currReceive, " id = ?  ", receive.ID)
	currReceive.PoNo = ""

	if removeItem {
		tx.Where("receive_id = ? ", currReceive.ID).Delete(dbmodels.ReceiveDetail{})
	}

	r := tx.Save(currReceive)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	tx.Commit()
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
