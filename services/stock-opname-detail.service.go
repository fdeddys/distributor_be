package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"time"
)

// StockOpnameDetailService ...
type StockOpnameDetailService struct {
}

// GetDataStockOpnameDetailPage ...
func (o StockOpnameDetailService) GetDataStockOpnameDetailPage(param dto.FilterStockOpname, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetStockOpnameDetailPage(param, offset, limit)

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
func (o StockOpnameDetailService) Save(stockMutationDetail *dbmodels.StockOpnameDetail) (errCode string, errDesc string) {

	if _, err := database.GetStockOpnameById(stockMutationDetail.StockOpnameID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SaveStockOpnameDetail(stockMutationDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func (o StockOpnameDetailService) SaveByUploadaData(stockOpnameId int64, stockUploads []dto.TemplateReportStockOpname) (errCode string, errDesc string) {

	total := 0
	hpp := 0
	stockOpname, _ := repository.GetStockOpnameById(stockOpnameId)
	for _, stockUpload := range stockUploads {

		curStock := int64(0)

		stock, errCode, _ := repository.GetStockByProductAndWarehouse(stockUpload.ProductID, stockOpname.WarehouseID)
		if errCode == constants.ERR_CODE_00 {
			curStock = stock.Qty
		}

		var stockOpnameDetail dbmodels.StockOpnameDetail
		stockOpnameDetail.LastUpdate = time.Now()
		stockOpnameDetail.LastUpdateBy = dto.CurrUser
		stockOpnameDetail.ProductID = stockUpload.ProductID
		stockOpnameDetail.Qty = stockUpload.Qty
		stockOpnameDetail.StockOpnameID = stockOpnameId
		stockOpnameDetail.UomID = stockUpload.UomID
		stockOpnameDetail.QtyOnSystem = curStock

		errcode, errdesc := database.SaveStockOpnameDetail(&stockOpnameDetail)
		if errcode != constants.ERR_CODE_00 {
			fmt.Println("Error ", stockUpload.ProductName, " ==? ", errdesc)
		}

		hpp = int(stockOpnameDetail.Hpp)
		total += ((int(stockOpnameDetail.Qty)) - int(curStock)) * hpp
	}
	fmt.Println("Update total Stock Opname ", total)
	stockOpname.Total = float32(total)
	database.SaveStockOpname(&stockOpname)
	fmt.Println("Finish add detail ")
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeleteStockOpnameDetailByID ...
func (o StockOpnameDetailService) DeleteStockOpnameDetailByID(stockMutationDetailId int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeleteStockOpnameDetailById(stockMutationDetailId); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateQty ...
func (o StockOpnameDetailService) UpdateQty(stockMutationDetailId, qty int64) (errCode string, errDesc string) {

	if err, errDesc := database.UpdateQtyReceiveSalesOrderDetail(stockMutationDetailId, qty); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
