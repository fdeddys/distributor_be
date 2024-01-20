package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/excel"
	"fmt"
	"time"
)

// PurchaseOrderService ...
type PurchaseOrderService struct {
}

// GetDataPage ...
func (r PurchaseOrderService) GetDataPage(param dto.FilterPurchaseOrder, page, limit, status int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetPurchaseOrderPage(param, offset, limit, status)

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

func (r PurchaseOrderService) ExportPurchaseOrder(param dto.FilterPurchaseOrder, status int) (bool, string) {
	res := false
	namaFile := "po_per.xlsx"
	data, _, err := database.GetPurchaseOrderPage(param, 1, 1000000, status)

	if err != nil {
		return res, ""
	}
	res = excel.ExportToExcelPo(data, namaFile)

	return res, namaFile
}

// GetDataPurchaseOrderByID ...
func (r PurchaseOrderService) GetDataPurchaseOrderByID(reveiveID int64) dbmodels.PurchaseOrder {

	var res dbmodels.PurchaseOrder
	// var err error
	res, _ = database.GetPurchaseOrderByPurchaseOrderID(reveiveID)

	return res
}

// Save ...
func (r PurchaseOrderService) Save(purchaseOrder *dbmodels.PurchaseOrder) (errCode, errDesc, purchaseOrderNo string, purchaseOrderID int64, status int8) {

	if purchaseOrder.ID == 0 {
		newNumber, errCode, errMsg := generateNewPurchaseOrderNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		purchaseOrder.PurchaserNo = newNumber
		purchaseOrder.Status = 10
	}
	purchaseOrder.LastUpdateBy = dto.CurrUser
	purchaseOrder.LastUpdate = time.Now()

	// fmt.Println("isi order ", order)
	err, errDesc, _, status := database.SavePurchaseOrder(purchaseOrder)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, purchaseOrder.PurchaserNo, purchaseOrder.ID, status
}

// ApprovePurchaseOrder ...
func (r PurchaseOrderService) ApprovePurchaseOrder(purchaseOrder *dbmodels.PurchaseOrder) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)
	err, errDesc := database.ApprovePurchaseOrder(purchaseOrder)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// RejectPurchaseOrder ...
// func (o OrderService) RejectPurchaseOrder(purchaseOrder *dbmodels.PurchaseOrder) (errCode, errDesc string) {

// 	// cek qty
// 	// validateQty()
// 	// fmt.Println("isi order ", order)
// 	err, errDesc := database.RejectPurchaseOrder(purchaseOrder)
// 	if err != constants.ERR_CODE_00 {
// 		return err, errDesc
// 	}
// 	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
// }

func generateNewPurchaseOrderNo() (newNumber string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "PO"

	err, number, errdesc := database.AddSequence(bln, thn, header)
	if err != constants.ERR_CODE_00 {
		return "", err, errdesc
	}
	newNumb := fmt.Sprintf("00000%v", number)
	// newNumb = newNumb[len(newNumb)-5 : len(newNumb)]
	// newNumber = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	fmt.Println("new numb bef : ", newNumb)
	runes := []rune(newNumb)
	newNumb = string(runes[len(newNumb)-5 : len(newNumb)])
	fmt.Println("new numb after : ", newNumb)

	// newNumb = newNumb[len(newNumb)-5 : len(newNumb)]
	newNumber = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newNumber, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// GetDataPurchaseOrderByID ...
func (r PurchaseOrderService) RejectPO(poID int64) dto.NoContentResponse {

	var res dto.NoContentResponse
	po, err := database.GetPurchaseOrderByPurchaseOrderID(poID)

	if err != nil {
		res.ErrCode = constants.ERR_CODE_40
		res.ErrDesc = err.Error()
		return res
	}
	if !(po.Status == 10 || po.Status == 20) {
		res.ErrCode = constants.ERR_CODE_99
		res.ErrDesc = "Status yang diizinkan cancel hanya outstanding "
		return res
	}

	code, desc := database.RejectPurchaseOrder(po)
	res.ErrCode = code
	res.ErrDesc = desc
	return res
}

func CalculateTotalPO(poID int64) {

	po, _ := database.GetPurchaseOrderByPurchaseOrderID(poID)
	totalPO := database.CountTotalPO(po.PurchaserNo)

	tax := float32(0)
	if po.Tax > 0 {
		tax = totalPO * po.Tax / 100
	}
	grandTotal := totalPO + tax

	database.UpdateTotal(po.ID, totalPO, tax, grandTotal)
}

// GetDataPurchaseOrderByID ...
func (r PurchaseOrderService) CancelSubmitPO(poID int64) dto.NoContentResponse {

	var res dto.NoContentResponse
	po, err := database.GetPurchaseOrderByPurchaseOrderID(poID)

	if err != nil {
		res.ErrCode = constants.ERR_CODE_40
		res.ErrDesc = err.Error()
		return res
	}
	if po.Status != 20 {
		res.ErrCode = constants.ERR_CODE_99
		res.ErrDesc = "Status yang diizinkan hanya PO yang sudah submit !"
		return res
	}

	code, desc := database.CancelSubmitPurchaseOrder(po)
	res.ErrCode = code
	res.ErrDesc = desc
	return res
}
