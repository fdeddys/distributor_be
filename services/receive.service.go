package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/utils/excel"
	"fmt"
	"time"
)

// ReceiveService ...
type ReceiveService struct {
}

// GetDataPage ...
func (r ReceiveService) GetDataPage(param dto.FilterReceive, page, limit, status int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReceivePage(param, offset, limit, status)

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

// GetDataReceiveByID ...
func (r ReceiveService) GetDataReceiveByID(reveiveID int64) dbmodels.Receive {

	var res dbmodels.Receive
	// var err error
	res, _ = database.GetReceiveByReceiveID(reveiveID)

	return res
}

// Save ...
func (r ReceiveService) Save(receive *dbmodels.Receive) (errCode, errDesc, receiveNo string, receiveID int64, status int8) {

	if receive.ID == 0 {
		newNumber, errCode, errMsg := generateNewReceiveNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		receive.ReceiveNo = newNumber
		receive.Status = 10
	}
	receive.LastUpdateBy = dto.CurrUser
	receive.LastUpdate = time.Now()

	// fmt.Println("isi order ", order)
	err, errDesc, _, status := database.SaveReceive(receive)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, receive.ReceiveNo, receive.ID, status
}

// Save ...
func (r ReceiveService) SaveByPO(receive *dbmodels.Receive) (errCode, errDesc, receiveNo string, receiveID int64, status int8) {

	var updateReceive dbmodels.Receive
	if receive.ID == 0 {
		newNumber, errCode, errMsg := generateNewReceiveNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		updateReceive.ReceiveNo = newNumber
		updateReceive.Status = 10
	} else {
		updateReceive, _ = database.GetReceiveByReceiveID(receive.ID)
	}
	updateReceive.LastUpdateBy = dto.CurrUser
	updateReceive.LastUpdate = time.Now()
	updateReceive.PoNo = receive.PoNo
	// fmt.Println("isi order ", order)
	err, errDesc, _, status := database.SaveReceive(&updateReceive)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}

	fmt.Println("Po No => ", receive.PoNo)
	if receive.PoNo != "" {
		// insertDetailFromPoNo
		poDetails := database.GetAllDataDetailPurchaseOrderByPoNo(receive.PoNo)

		for _, poDetail := range poDetails {

			product, _, _ := database.FindProductByID(poDetail.ProductID)

			var receiveDetail dbmodels.ReceiveDetail
			receiveDetail.Disc1 = poDetail.Disc1
			receiveDetail.Disc2 = poDetail.Disc2
			receiveDetail.Hpp = product.Hpp
			receiveDetail.LastUpdate = time.Now()
			receiveDetail.LastUpdateBy = dto.CurrUser
			receiveDetail.Price = poDetail.Price
			receiveDetail.ProductID = poDetail.ProductID
			receiveDetail.Qty = poDetail.Qty
			receiveDetail.ReceiveID = receive.ID
			receiveDetail.UomID = poDetail.UomID
			database.SaveReceiveDetail(&receiveDetail)
		}
		database.UpdatePoPaid(receive.PoNo)
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, receive.ReceiveNo, receive.ID, status
}

// ApproveReceive ...
func (r ReceiveService) ApproveReceive(order *dbmodels.Receive) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)
	err, errDesc := database.SaveReceiveApprove(order)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// RejectReceive ...
func (o ReceiveService) RejectReceive(receive *dbmodels.Receive) (errCode, errDesc string) {

	// cek qty
	// validateQty()
	// fmt.Println("isi order ", order)
	err, errDesc := database.RejectReceive(receive)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func generateNewReceiveNo() (newNumber string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "RV"

	err, number, errdesc := database.AddSequence(bln, thn, header)
	if err != constants.ERR_CODE_00 {
		return "", err, errdesc
	}
	newNumb := fmt.Sprintf("00000%v", number)
	newNumb = newNumb[len(newNumb)-5 : len(newNumb)]
	newNumber = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newNumber, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// ApproveReceive ...
func (r ReceiveService) RemovePO(order *dbmodels.Receive, removeItem bool) (errCode, errDesc string) {

	// fmt.Println("isi order ", order)
	err, errDesc := database.RemovePO(order, removeItem)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func (r ReceiveService) ExportReceive(param dto.FilterReceive, status int) (bool, string) {
	res := false
	namaFile := "receive_per.xlsx"
	data, _, err := database.GetReceivePage(param, 0, 1000000, status)

	if err != nil {
		return res, ""
	}
	res = excel.ExportToExcelReceive(data, namaFile)

	return res, namaFile
}
