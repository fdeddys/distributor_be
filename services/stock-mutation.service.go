package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"time"
)

// StockMutationService ...
type StockMutationService struct {
}

// GetDataStockMutationById ...
func (o StockMutationService) GetDataStockMutationById(stockMutationID int64) dbmodels.StockMutation {

	var res dbmodels.StockMutation
	// var err error
	res, _ = database.GetStockMutationById(stockMutationID)

	return res
}

// GetDataPage ...
func (o StockMutationService) GetDataPage(param dto.FilterStockMutation, page int, limit int, internalStatus int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetStockMutationPage(param, offset, limit, internalStatus)

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
func (o StockMutationService) Save(stockMutation *dbmodels.StockMutation) (errCode, errDesc, stockMutationNo string, stockMutationID int64, status int8) {

	if stockMutation.ID == 0 {
		newStockMutationNo, errCode, errMsg := generateNewStockMutationNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		stockMutation.StockMutationNo = newStockMutationNo
		stockMutation.Status = 10
	}
	stockMutation.LastUpdateBy = dto.CurrUser
	stockMutation.LastUpdate = time.Now()

	err, errDesc, newID, status := database.SaveStockMutation(stockMutation)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, stockMutation.StockMutationNo, newID, status
}

// Approve ...
func (o StockMutationService) Approve(stockMutation *dbmodels.StockMutation) (errCode, errDesc string) {

	// cek qty
	valid, errCode, errDesc := validateMutation(stockMutation.ID, stockMutation.WarehouseSourceID)
	if !valid {
		return errCode, errDesc
	}
	err, errDesc := database.SaveStockMutationApprove(stockMutation)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func validateMutation(stockMutationID, warehouseSourceID int64) (isValid bool, errCode, errDesc string) {

	pesan := ""

	stockMutationDetails := database.GetAllDataStockMutationDetail(stockMutationID)
	for idx, stockMutationDetail := range stockMutationDetails {
		fmt.Println("idx -> ", idx)

		product, errCodeProd, _ := database.FindProductByID(stockMutationDetail.ProductID)
		if errCodeProd != constants.ERR_CODE_00 {
			pesan += fmt.Sprintf(" [%v] Product not found or inactive !", stockMutationDetail.ProductID)
			break
		}

		checkStock, errcode, errDesc := database.GetStockByProductAndWarehouse(product.ID, warehouseSourceID)
		if errcode != constants.ERR_CODE_00 {
			pesan += fmt.Sprintf(" [%v] %v", product.Name, errDesc)
			break
		}
		curQty := checkStock.Qty
		stockMutationQty := stockMutationDetail.Qty

		if stockMutationQty > curQty {
			pesan += fmt.Sprintf(" [%v] qty stock Mutation = %v more than qty stock = %v!", product.Name, stockMutationQty, curQty)
			break
		}
	}
	if pesan == "" {
		return true, "", ""
	}
	return false, constants.ERR_CODE_80, pesan
}

// Reject ...
func (o StockMutationService) Reject(stockMutation *dbmodels.StockMutation) (errCode, errDesc string) {

	err, errDesc := database.RejectStockMutation(stockMutation)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func generateNewStockMutationNo() (newStockMutationNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "MT"

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
	newStockMutationNo = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newStockMutationNo, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}
