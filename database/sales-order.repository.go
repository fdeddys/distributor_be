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

//SaveSalesOrderNo ...
func SaveSalesOrderNo(order *dbmodels.SalesOrder) (errCode string, errDesc string, id int64, status int8) {

	fmt.Println(" Update Sales Order numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	r := db.Save(&order)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	// fmt.Println("Order [database]=> order id", order.OrderNo)

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, order.ID, order.Status
}

//RejectSalesOrder ...
func RejectSalesOrder(order *dbmodels.SalesOrder) (errCode string, errDesc string) {

	fmt.Println(" Reject Sales Order numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	// salesOrder, _ := GetOrderByOrderNo(order.SalesOrderNo)
	// salesOrderDetails := GetAllDataDetail(order.ID)
	// for idx, orderDetail := range salesOrderDetails {
	// 	fmt.Println("idx -> ", idx)

	// 	updateStockInsertHistoryReject(salesOrder, orderDetail)

	// }

	r := db.Model(&dbmodels.SalesOrder{}).Where("id =?", order.ID).Update(dbmodels.SalesOrder{Status: 30})
	// r := db.Save(&order)
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// SaveSalesOrderApprove ...
func SaveSalesOrderApprove(order *dbmodels.SalesOrder) (errCode string, errDesc string) {

	fmt.Println(" Approve Sales Order ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// update stock untuk penjualan tunai dan kredit
	// update stock -> jika penjualan kredit , klo tunai update saat payment
	// update history stock
	// hitung ulang
	var total float32
	var grandTotal float32
	total = 0
	grandTotal = 0
	salesOrderDetails := GetAllDataDetail(order.ID)
	for idx, orderDetail := range salesOrderDetails {
		fmt.Println("idx -> ", idx)

		// if order.IsCash == false {
		updateStockInsertHistory(*order, orderDetail)
		// }
		disc := orderDetail.Price * float32(orderDetail.QtyOrder) * orderDetail.Disc1 / 100
		total = total + (orderDetail.Price*float32(orderDetail.QtyOrder) - disc)
	}

	db.Debug().LogMode(true)
	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	grandTotal = total
	if order.Tax != 0 {
		grandTotal = total + total*order.Tax/100
	}
	order.GrandTotal = grandTotal
	order.Total = total
	order.SalesmanID = dto.CurrUserId
	order.LastUpdateBy = dto.CurrUser
	order.LastUpdate = time.Now()
	order.Status = 20
	r := db.Save(&order)
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

func updateStockInsertHistory(order dbmodels.SalesOrder, orderDetail dbmodels.SalesOrderDetail) {
	product, _, _ := FindProductByID(orderDetail.ProductID)
	// if errCodeProd != constants.ERR_CODE_00 {
	// 	tx.Rollback()
	// 	return errCodeProd, errDescProd
	// }
	// curQty := product.QtyStock
	checkStock, _, _ := GetStockByProductAndWarehouse(product.ID, order.WarehouseID)
	curQty := checkStock.Qty

	updateQty := curQty - orderDetail.QtyOrder

	fmt.Println("cur qty =", curQty, " update =", updateQty)
	var historyStock dbmodels.HistoryStock
	historyStock.Code = product.Code
	if order.IsCash {
		historyStock.Description = "Direct Sales"
	} else {
		historyStock.Description = "Sales Order"
	}
	historyStock.Hpp = product.Hpp
	historyStock.WarehouseID = order.WarehouseID
	historyStock.Name = product.Name
	historyStock.Price = orderDetail.Price
	historyStock.ReffNo = order.SalesOrderNo
	historyStock.TransDate = order.OrderDate
	historyStock.Debet = 0
	historyStock.Kredit = orderDetail.QtyOrder
	historyStock.Saldo = updateQty
	historyStock.Disc1 = orderDetail.Disc1 * orderDetail.Price / 100
	historyStock.LastUpdate = time.Now()
	historyStock.LastUpdateBy = dto.CurrUser

	UpdateStockProductByID(orderDetail.ProductID, updateQty, order.WarehouseID)
	SaveHistory(historyStock)
}

// GetOrderByOrderNo ...
func GetOrderByOrderNo(orderNo string) (dbmodels.SalesOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	order := dbmodels.SalesOrder{}

	err := db.Preload("Customer").Where(" sales_order_no = ?  ", orderNo).First(&order).Error

	return order, err

}

// GetSalesOrderByOrderId ...
func GetSalesOrderByOrderId(orderID int64) (dbmodels.SalesOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	order := dbmodels.SalesOrder{}

	err := db.Preload("Customer").Preload("Salesman").Where(" id = ?  ", orderID).First(&order).Error

	return order, err

}

// GetOrderPage ...
func GetOrderPage(param dto.FilterOrder, offset, limit, internalStatus int) ([]dbmodels.SalesOrder, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var orders []dbmodels.SalesOrder
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&orders).Error
		if err != nil {
			return orders, 0, err
		}
		return orders, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysOrders(db, offset, limit, internalStatus, &orders, param, errQuery)
	go AsyncQueryCountsOrders(db, &total, internalStatus, &orders, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return orders, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return orders, 0, resErrCount
	}
	return orders, total, nil
}

func getParam(param dto.FilterOrder, status int) (merchantCode, orderNumber string, byStatus, isCash bool) {

	merchantCode = "%"

	orderNumber = param.OrderNumber
	if orderNumber == "" {
		orderNumber = "%"
	} else {
		orderNumber = "%" + param.OrderNumber + "%"
	}

	byStatus = true
	if status == -1 {
		byStatus = false
	}

	isCash = param.IsCash
	return
}

// AsyncQueryCountsOrders ...
func AsyncQueryCountsOrders(db *gorm.DB, total *int, status int, orders *[]dbmodels.SalesOrder, param dto.FilterOrder, resChan chan error) {

	merchantCode, orderNumber, byStatus, isCash := getParam(param, status)

	fmt.Println("ISI MERCHANT ", merchantCode, " orderNumber ", orderNumber, "  status ", status, " fill status ", byStatus, "  Is Cash : ", isCash)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&orders).Where(" ( (status = ?) or ( not ?) ) AND  COALESCE(sales_order_no, '') ilike ? AND order_date between ? and ?  and is_cash = ? ", status, byStatus, orderNumber, param.StartDate, param.EndDate, isCash).Count(&*total).Error
	} else {
		err = db.Model(&orders).Where(" ( (status = ?) or ( not ?) ) AND COALESCE(sales_order_no,'') ilike ? and is_cash = ?", status, byStatus, orderNumber, isCash).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysOrders ...
func AsyncQuerysOrders(db *gorm.DB, offset int, limit int, status int, orders *[]dbmodels.SalesOrder, param dto.FilterOrder, resChan chan error) {

	var err error

	merchantCode, orderNumber, byStatus, isCash := getParam(param, status)

	fmt.Println("ISI MERCHANT ", merchantCode, " order no ", orderNumber, "  status ", status, " fill status ", byStatus)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("Customer").Preload("Salesman").Order("id DESC").Offset(offset).Limit(limit).Find(&orders, " ( ( status = ?) or ( not ?) ) AND COALESCE(sales_order_no, '') ilike ? AND order_date between ? and ?  and is_cash = ? ", status, byStatus, orderNumber, param.StartDate, param.EndDate, isCash).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("Customer").Preload("Salesman").Order("id DESC").Find(&orders, " ( ( status = ?) or ( not ?) ) AND COALESCE(sales_order_no,'') ilike ?  and is_cash = ? ", status, byStatus, orderNumber, isCash).Error
		if err != nil {
			fmt.Println("error --> ", err)
		}
		fmt.Println("order--> ", orders)

	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

func GetSalesOrderForPayment(param dto.FilterOrder, offset, limit int) ([]dbmodels.SalesOrder, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var orders []dbmodels.SalesOrder
	var total int

	total = 0
	err := db.Offset(offset).Limit(limit).Preload("Customer").Preload("Salesman").Order("order_date DESC").Find(&orders, "  ( status  in ('20','40') ) AND (is_paid is null or is_paid = false) AND customer_id = ? ", param.CustomerID).Error

	if err != nil {
		return orders, total, err
	}

	var cekOorders []dbmodels.SalesOrder
	db.Model(&cekOorders).Where("  ( status  in  ('20','40')  ) AND (is_paid is null or is_paid = false) AND customer_id = ? ", param.CustomerID).Count(&total)

	return orders, total, nil
}
