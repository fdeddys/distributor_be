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

// ReturnReceiveService ...
type ReturnReceiveService struct {
}

// GetDataOrderById ...
func (o ReturnReceiveService) GetDataReturnById(orderID int64) dbmodels.ReturnReceive {

	var res dbmodels.ReturnReceive
	// var err error
	res, _ = database.GetReturnReceiveById(orderID)

	return res
}

// GetDataReceiveReturnPage ...
func (o ReturnReceiveService) GetDataReturnReceivePage(param dto.FilterReturnReceive, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReturnReceivePage(param, offset, limit)

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
func (o ReturnReceiveService) Save(receiveReturn *dbmodels.ReturnReceive) (errCode, errDesc, returnReceiveNo string, returnReceiveID int64, status int8) {

	if receiveReturn.ID == 0 {
		newReceiveReturnNo, errCode, errMsg := generateReceiveReturnNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		receiveReturn.ReturnReceiveNo = newReceiveReturnNo
		receiveReturn.Status = 10
	}
	receiveReturn.LastUpdateBy = dto.CurrUser
	receiveReturn.LastUpdate = time.Now()

	// fmt.Println("isi order ", order)
	err, errDesc, newID, status := database.SaveReturnReceive(receiveReturn)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, receiveReturn.ReturnReceiveNo, newID, status
}

func calculateTotalReturnReceive(returnReceiveID int64){
	
	fmt.Println("Calculate Total receive....")
	db := database.GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("defar, roll back ")
			tx.Rollback()
		}
	}()

	var grandTotal float32
	grandTotal=0

	returnReceive,_ :=database.GetReceiveByReceiveID(returnReceiveID)
	returnReceiveDetails := database.GetAllReceiveReturnDetail(returnReceiveID)

	var total float32
	for _, returnReceiveDetail := range returnReceiveDetails {

		total = total + (returnReceiveDetail.Price * float32(returnReceiveDetail.Qty))

	}
	grandTotal = total

	if returnReceive.Tax != 0 {
		grandTotal = total * 1.1
	}

	tx.Model(&dbmodels.ReturnReceive{}).
		Where("id = ?", returnReceiveID).
		Update(dbmodels.ReturnReceive{
			GrandTotal:   grandTotal,
			Total:        total,
			LastUpdateBy: dto.CurrUser,
			LastUpdate:   util.GetCurrDate(),
		})

	tx.Commit()

}

func generateReceiveReturnNo() (newOrderNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "RR"

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
func (o ReturnReceiveService) Approve(returnReceiveID int64) (errCode, errDesc string) {

	db := database.GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("defar, roll back ")
			tx.Rollback()
		}
	}()

	// cek qty
	returnReceive, err := database.GetReturnReceiveById(returnReceiveID)
	if err != nil {
		tx.Rollback()
		return errCode, errDesc
	}

	returnReceiveDetails := database.GetAllReceiveReturnDetail(returnReceiveID)

	var total float32
	for _, returnReceiveDetail := range returnReceiveDetails {

		product, _, _ := database.FindProductByID(returnReceiveDetail.ProductID)
		stock, _, _ := database.GetStockByProductAndWarehouse(returnReceiveDetail.ProductID, returnReceive.WarehouseID)

		var newStock int64
		curStock := stock.Qty
		newStock = curStock - returnReceiveDetail.Qty
		fmt.Println("cur Stock = " , curStock)
		fmt.Println("new Stock = " , newStock)

		// db.Model(&stock).Select("Qty","LastUpdateBy","LastUpdate").Updates(map[string]interface{}{"Qty": curStock - returnReceiveDetail.Qty, "LastUpdateBy": dto.CurrUser, "LastUpdate": util.GetCurrDate()})
		

		// tx.Model(&dbmodels.Stock{}).
		// Where("id = ?", stock.ID).
		// Select("Qty","LastUpdateBy","LastUpdate").
		// Update(dbmodels.Stock{
		// 	Qty: newStock,
		// 	LastUpdateBy: dto.CurrUser,
		// 	LastUpdate:   util.GetCurrDate(),
		// })
		
		
		stock.LastUpdateBy = dto.CurrUser
		stock.LastUpdate = util.GetCurrDate()
		stock.Qty = newStock
		db.Save(&stock)
		

		// TODO hitung disc
		// hpp := (float32(curStock * stock.Qty)) + (float32(returnReceiveDetail.Qty)*returnReceiveDetail.Price)/float32(stock.Qty+returnReceiveDetail.Qty)
		hpp := product.Hpp

		var history dbmodels.HistoryStock
		history.Code = product.Code
		history.WarehouseID = returnReceive.WarehouseID
		history.Debet = 0
		history.Description = "Return Receive"
		history.Hpp = hpp
		history.Kredit = returnReceiveDetail.Qty
		history.LastUpdate = util.GetCurrDate()
		history.LastUpdateBy = dto.CurrUser
		history.Name = product.Name
		history.ReffNo = returnReceive.ReturnReceiveNo
		history.Price = returnReceiveDetail.Price
		history.Saldo = newStock
		history.TransDate = returnReceive.ReturnReceiveDate
		tx.Save(&history)

		total = total + (returnReceiveDetail.Price * float32(returnReceiveDetail.Qty))

	}

	var grandTotal float32
	// total = 0
	grandTotal = total

	if returnReceive.Tax != 0 {
		grandTotal = total * 1.1
	}

	tx.Model(&dbmodels.ReturnReceive{}).
		Where("id = ?", returnReceive.ID).
		Update(dbmodels.ReturnReceive{
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
func (o ReturnReceiveService) Reject(returnOrder *dbmodels.ReturnReceive) (errCode, errDesc string) {

	err, errDesc := database.RejectReturnReceive(returnOrder.ID)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
