package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// StockOpnameService ...
type StockOpnameService struct {
}

// GetDataStockOpnameById ...
func (o StockOpnameService) GetDataStockOpnameById(stockOpnameID int64) dbmodels.StockOpname {

	var res dbmodels.StockOpname
	// var err error
	res, _ = database.GetStockOpnameById(stockOpnameID)

	return res
}

// GetDataPage ...
func (o StockOpnameService) GetDataPage(param dto.FilterStockOpname, page int, limit int, internalStatus int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetStockOpnamePage(param, offset, limit, internalStatus)

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
func (o StockOpnameService) Save(stockOpname *dbmodels.StockOpname) (errCode, errDesc, stockOpnameNo string, stockOpnameID int64, status int8) {

	if stockOpname.ID == 0 {
		newStockOpnameNo, errCode, errMsg := generateNewStockOpnameNo()
		if errCode != constants.ERR_CODE_00 {
			return errCode, errMsg, "", 0, 0
		}
		stockOpname.StockOpnameNo = newStockOpnameNo
		stockOpname.Status = 10
	}
	stockOpname.LastUpdateBy = dto.CurrUser
	stockOpname.LastUpdate = time.Now()

	err, errDesc, newID, status := database.SaveStockOpname(stockOpname)
	if err != constants.ERR_CODE_00 {
		return err, errDesc, "", 0, 0
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, stockOpname.StockOpnameNo, newID, status
}

// Approve ...
func (o StockOpnameService) Approve(stockOpname *dbmodels.StockOpname) (errCode, errDesc string) {

	err, errDesc := database.SaveStockOpnameApprove(stockOpname)
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func generateNewStockOpnameNo() (newStockOpnameNo string, errCode string, errMsg string) {

	t := time.Now()
	bln := t.Format("01")
	thn := t.Format("06")
	header := "OP"

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
	newStockOpnameNo = fmt.Sprintf("%v%v%v%v", header, thn, bln, newNumb)

	return newStockOpnameNo, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG

}

// Approve ...
func (o StockOpnameService) GenerateTemplateStockOpname(warehouseID int64) {

	datas := generateData(warehouseID)
	exportToCSV(datas)
}

func exportToCSV(datas []dto.TemplateReportStockOpname) {
	filename := "stock-opname.csv"
	csvFile, _ := os.Create(filename)

	defer func() {
		csvFile.Close()
	}()

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ';'
	result := util.ToSliceData(datas)

	for _, data := range result {
		csvwriter.Write(data)
	}
	csvwriter.Flush()
	csvFile.Close()

}

func generateData(warehouseID int64) []dto.TemplateReportStockOpname {

	datas := database.FindDataStockByWarehouseID(warehouseID)
	return datas
}

// Approve ...
func (o StockOpnameService) RecalculateTotal() (errCode, errDesc string) {

	err, errDesc := database.RecalculateTotal()
	if err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
