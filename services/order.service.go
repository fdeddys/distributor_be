package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"time"
)

// OrderService ...
type OrderService struct {
}

// GetDataOrderById ...
func (o OrderService) GetDataOrderById(orderID int64) dbmodels.SalesOrder {

	var res dbmodels.SalesOrder
	// var err error
	res, _ = database.GetSalesOrderByOrderId(orderID)

	return res
}

// GetDataOrderById ...
func (o OrderService) GetTotalOrderById(orderID int64) float32 {

	// var err error
	grandTotal := float32(0)
	totalDetails := database.GetAllDataDetail(orderID)

	for _, detail := range totalDetails {
		total := float32(0)

		total = detail.Price * float32(detail.QtyOrder)
		total = total - (total * detail.Disc1 / 100)
		total = total - (total * detail.Disc2 / 100)

		grandTotal += total
	}

	return grandTotal
}

// GetDataPage ...
func (o OrderService) GetDataPage(param dto.FilterOrder, page int, limit int, internalStatus int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetOrderPage(param, offset, limit, internalStatus)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = limit

	return res
}

// Save ...
func (o OrderService) Save(order *dbmodels.SalesOrder) (errCode, errDesc, orderNo string, orderID int64, status int8) {

	if order.ID == 0 {
		newOrderNo, errCode, errMsg := generateNewOrderNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		order.SalesOrderNo = newOrderNo
		order.Status = 10
		order.SalesmanID = dto.CurrUserId
	}
	order.LastUpdateBy = dto.CurrUser
	order.LastUpdate = time.Now()

	// fmt.Println("isi order ", order)
	err, errDesc, newID, status := database.SaveSalesOrderNo(order)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, order.SalesOrderNo, newID, status
}

// Approve ...
func (o OrderService) Approve(order *dbmodels.SalesOrder) (errCode, errDesc string) {

	// cek qty
	valid, errCode, errDesc := validateQty(order.ID, order.WarehouseID)
	if !valid {
		return errCode, errDesc
	}
	// fmt.Println("isi order ", order)
	err, errDesc := database.SaveSalesOrderApprove(order)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func validateQty(orderID, warehouseID int64) (isValid bool, errCode, errDesc string) {

	pesan := ""

	salesOrderDetails := database.GetAllDataDetail(orderID)
	for idx, orderDetail := range salesOrderDetails {
		fmt.Println("idx -> ", idx)

		product, errCodeProd, _ := database.FindProductByID(orderDetail.ProductID)
		if errCodeProd != constants.ERR_CODE_00 {
			pesan += fmt.Sprintf(" [%v] Product not found or inactive !", orderDetail.ProductID)
			break
		}

		checkStock, errcode, errDesc := database.GetStockByProductAndWarehouse(product.ID, warehouseID)
		if errcode != constants.ERR_CODE_00 {
			pesan += fmt.Sprintf(" [%v] %v", product.Name, errDesc)
			break
		}
		curQty := checkStock.Qty
		orderQty := orderDetail.QtyOrder

		if orderQty > curQty {
			pesan += fmt.Sprintf(" [%v] qty order = %v more than qty stock = %v!", product.Name, orderQty, curQty)
			break
		}
	}
	if pesan == "" {
		return true, "", ""
	}
	return false, constants.ERR_CODE_80, pesan
}

// Reject ...
func (o OrderService) Reject(order *dbmodels.SalesOrder) (errCode, errDesc string) {

	// cek qty
	// validateQty()
	// fmt.Println("isi order ", order)
	err, errDesc := database.RejectSalesOrder(order)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func generateNewOrderNo() (newOrderNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "SO"

	err, number, errdesc := database.AddSequence(bln, thn, header)
	if err != constants.ERR_CODE_00 {
		return "", err, errdesc
	}
	newNumb := fmt.Sprintf("00000%v", number)
	fmt.Println("new numb bef : ", newNumb)
	runes := []rune(newNumb)
	newNumb = string(runes[len(newNumb)-5 : len(newNumb)])
	fmt.Println("new numb after : ", newNumb)

	// newNumb = newNumb[len(newNumb)-5 : len(newNumb)]
	newOrderNo = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newOrderNo, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// // PrintPdf ...
// func (o OrderService) PrintPdf(order *dbmodels.Order) (errCode string, errDesc string) {

// 	// if err, errDesc := database.SaveSalesOrderNo(order); err != constants.ERR_CODE_00 {
// 	// 	return err, errDesc
// 	// }

// 	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
// }

// CreateInvoice ...
func (o OrderService) CreateInvoice(orderID int64) (errCode, errDesc string) {

	db := database.GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// cek qty
	salesOrder, err := database.GetSalesOrderByOrderId(orderID)
	if err != nil {
		return errCode, errDesc
	}

	salesOrderDetails := database.GetAllDataDetail(salesOrder.ID)

	var total float32
	for _, sod := range salesOrderDetails {

		product, _, _ := database.FindProductByID(sod.ProductID)
		stock, _, _ := database.GetStockByProductAndWarehouse(sod.ProductID, salesOrder.WarehouseID)

		curStock := stock.Qty
		newStock := curStock - sod.QtyReceive

		db.Model(&dbmodels.Stock{}).
			Where("id = ?", stock.ID).
			Update(dbmodels.Stock{
				Qty:          newStock,
				LastUpdateBy: dto.CurrUser,
				LastUpdate:   util.GetCurrDate(),
			})

		var history dbmodels.HistoryStock
		history.Code = product.Code
		history.Debet = 0
		history.Description = "Sales Order"
		history.Hpp = product.Hpp
		history.Kredit = sod.QtyReceive
		history.LastUpdate = util.GetCurrDate()
		history.LastUpdateBy = dto.CurrUser
		history.Name = product.Name
		history.ReffNo = salesOrder.SalesOrderNo
		history.Price = sod.Price
		history.Saldo = newStock
		history.TransDate = salesOrder.OrderDate
		tx.Save(&history)

		total = total + (sod.Price * float32(sod.QtyOrder))

	}

	var grandTotal float32
	invNo, _, _ := generateNewInvoiceNo()
	total = 0
	grandTotal = total

	if salesOrder.Tax != 0 {
		grandTotal = total * 1.1
	}
	// salesOrder.GrandTotal = grandTotal
	// salesOrder.Status = 40
	// salesOrder.LastUpdate = util.GetCurrDate()
	// salesOrder.LastUpdateBy = dto.CurrUser
	// salesOrder.InvoiceNo = invNo
	// salesOrder.Total = total
	// db.Save(&salesOrder)

	tx.Model(&dbmodels.SalesOrder{}).
		Where("id = ?", salesOrder.ID).
		Update(dbmodels.SalesOrder{
			GrandTotal:   grandTotal,
			Status:       40,
			InvoiceNo:    invNo,
			Total:        total,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		})

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func generateNewInvoiceNo() (newOrderNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "IV"

	err, number, errdesc := database.AddSequence(bln, thn, header)
	if err != constants.ERR_CODE_00 {
		return "", err, errdesc
	}
	newNumb := fmt.Sprintf("00000%v", number)
	fmt.Println("new numb bef : ", newNumb)
	runes := []rune(newNumb)
	newNumb = string(runes[len(newNumb)-5 : len(newNumb)])
	fmt.Println("new numb after : ", newNumb)

	// newNumb = newNumb[len(newNumb)-5 : len(newNumb)]
	newOrderNo = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newOrderNo, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// GetDataForSalesOrderPage ...
func (o OrderService) GetDataForSalesOrderPage(param dto.FilterOrder, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetSalesOrderForPayment(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = limit

	return res
}
