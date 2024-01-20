package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"time"
)

// ReturnSalesOrderService ...
type ReturnSalesOrderService struct {
}

// GetDataOrderById ...
func (o ReturnSalesOrderService) GetDataOrderReturnById(orderID int64) dbmodels.ReturnSalesOrder {

	var res dbmodels.ReturnSalesOrder
	// var err error
	res, _ = database.GetSalesOrderReturnById(orderID)

	return res
}

// GetDataSalesOrderReturnPage ...
func (o ReturnSalesOrderService) GetDataSalesOrderReturnPage(param dto.FilterOrderReturnDetail, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReturnSalesOrderPage(param, offset, limit)

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
func (o ReturnSalesOrderService) Save(salesOrderReturn *dbmodels.ReturnSalesOrder) (errCode, errDesc, salesOrderReturnNo string, salesOrderReturnID int64, status int8) {

	if salesOrderReturn.ID == 0 {
		newSalesOrderReturnNo, errCode, errMsg := generateSalesOrderReturnNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		salesOrderReturn.ReturnSalesOrderNo = newSalesOrderReturnNo
		salesOrderReturn.Status = 10
		salesOrderReturn.SalesmanID = salesOrderReturn.SalesmanID
	}
	salesOrderReturn.LastUpdateBy = dto.CurrUser
	salesOrderReturn.LastUpdate = time.Now()

	// fmt.Println("isi order ", order)
	err, errDesc, newID, status := database.SaveSalesOrderReturn(salesOrderReturn)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, salesOrderReturn.ReturnSalesOrderNo, newID, status
}

func generateSalesOrderReturnNo() (newOrderNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "SR"

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

// CreateInvoice ...
func (o ReturnSalesOrderService) Approve(salesOrderReturnID int64) (errCode, errDesc string) {

	db := database.GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// cek qty
	salesOrderReturn, err := database.GetSalesOrderReturnById(salesOrderReturnID)
	if err != nil {
		return errCode, errDesc
	}

	salesOrderReturnDetails := database.GetAllSalesOrderReturnDetail(salesOrderReturnID)

	var total float32
	for _, sord := range salesOrderReturnDetails {

		product, _, _ := database.FindProductByID(sord.ProductID)
		stock, _, _ := database.GetStockByProductAndWarehouse(sord.ProductID, salesOrderReturn.WarehouseID)

		curStock := stock.Qty
		newStock := curStock + sord.Qty

		db.Model(&dbmodels.Stock{}).
			Where("id = ?", stock.ID).
			Update(dbmodels.Stock{
				Qty:          newStock,
				LastUpdateBy: dto.CurrUser,
				LastUpdate:   util.GetCurrDate(),
			})

		// TODO hitung disc
		hpp := (float32(int64(product.Hpp) * curStock)) + (float32(sord.Qty)*sord.Price)/float32(curStock+sord.Qty)
		var history dbmodels.HistoryStock
		history.Code = product.Code
		history.Debet = sord.Qty
		history.Description = "Sales Order Return"
		history.Hpp = hpp
		history.Kredit = 0
		history.LastUpdate = util.GetCurrDate()
		history.LastUpdateBy = dto.CurrUser
		history.Name = product.Name
		history.ReffNo = salesOrderReturn.ReturnSalesOrderNo
		history.Price = sord.Price
		history.Saldo = newStock
		history.TransDate = salesOrderReturn.ReturnSalesOrderDate
		history.WarehouseID = salesOrderReturn.WarehouseID
		history.Total = salesOrderReturn.Total
		db.Save(&history)

		total = total + (sord.Price * float32(sord.Qty))

	}

	var grandTotal float32
	// total = 0
	grandTotal = total

	fmt.Println("total = ", total, "  grand total = ", grandTotal)

	if salesOrderReturn.Tax != 0 {
		grandTotal = total * 1.1
	}

	db.Model(&dbmodels.ReturnSalesOrder{}).
		Where("id = ?", salesOrderReturn.ID).
		Update(dbmodels.ReturnSalesOrder{
			GrandTotal:   grandTotal,
			Status:       20,
			Total:        total,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		})

	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// Reject ...
func (o ReturnSalesOrderService) Reject(returnOrder *dbmodels.ReturnSalesOrder) (errCode, errDesc string) {

	err, errDesc := database.RejectReturnSalesOrder(returnOrder.ID)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetDataForSalesOrderPage ...
func (o ReturnSalesOrderService) GetDataForSalesOrderReturnPage(param dto.FilterOrderReturnDetail, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetSalesOrderReturnForPayment(param, offset, limit)

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
