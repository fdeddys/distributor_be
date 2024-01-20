package controllers

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/services"
	"distribution-system-be/utils/util"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// StockOpnameController ...
type StockOpnameController struct {
	DB *gorm.DB
}

// StockOpnameService ...
var StockOpnameService = new(services.StockOpnameService)

// GetByStockOpnameId ...
func (s *StockOpnameController) GetByStockOpnameById(c *gin.Context) {

	res := dbmodels.StockOpname{}

	stockMutationID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = StockOpnameService.GetDataStockOpnameById(stockMutationID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// GetStockOpnamePage ...
func (s *StockOpnameController) GetStockOpnamePage(c *gin.Context) {
	req := dto.FilterStockOpname{}
	res := models.ResponsePagination{}

	page, errPage := strconv.Atoi(c.Param("page"))
	if errPage != nil {
		logs.Info("error", errPage)
		res.Error = errPage.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	count, errCount := strconv.Atoi(c.Param("count"))
	if errCount != nil {
		logs.Info("error", errPage)
		res.Error = errCount.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("page->", page, "count->", count)

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	temp, _ := json.Marshal(req)
	log.Println("searchName-->", string(temp))
	log.Println("is release ", req.InternalStatus)

	intStatus := -1
	if intVal, errconv := strconv.Atoi(req.InternalStatus); errconv == nil {
		intStatus = intVal
	}
	res = StockOpnameService.GetDataPage(req, page, count, intStatus)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *StockOpnameController) Save(c *gin.Context) {

	req := dbmodels.StockOpname{}
	body := c.Request.Body
	res := dto.StockOpnameSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales StockOpname stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, stockMutationNo, stockMutationID, status := StockOpnameService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.StockOpnameNo = stockMutationNo
	res.ID = stockMutationID
	res.Status = status
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *StockOpnameController) Approve(c *gin.Context) {

	req := dbmodels.StockOpname{}
	body := c.Request.Body
	res := dto.StockOpnameSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales StockOpname stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		return
	}

	errCode, errMsg := StockOpnameService.Approve(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)
	return
}

// PrintSO ...
// func (s *StockOpnameController) PrintStockOpnameForm(c *gin.Context) {

// 	stockOpnameID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if errPage != nil {
// 		logs.Info("error", errPage)
// 		c.JSON(http.StatusBadRequest, "id not supplied")
// 		c.Abort()
// 		return
// 	}

// 	report.GenerateStockOpnameReport(stockOpnameID)

// 	header := c.Writer.Header()
// 	header["Content-type"] = []string{"application/x-pdf"}
// 	header["Content-Disposition"] = []string{"attachment; filename= tes.pdf"}

// 	file, _ := os.Open("stock-mutation.pdf")

// 	io.Copy(c.Writer, file)
// 	return
// }

// DownloadTemplate ...
func (s *StockOpnameController) DownloadTemplate(c *gin.Context) {
	// fmt.Println("-------->", req)

	warehouseID, err := strconv.ParseInt(c.Param("warehouseId"), 10, 64)
	if err != nil {
		fmt.Println("Erro => ", err)
		return
	}

	StockOpnameService.GenerateTemplateStockOpname(warehouseID)

	fmt.Println("download template")
	header := c.Writer.Header()
	header["Content-type"] = []string{"text/csv"}
	header["Content-Disposition"] = []string{"attachment; filename=stock-opname.csv"}

	file, _ := os.Open("stock-opname.csv")

	io.Copy(c.Writer, file)
	return
}

// UploadTemplate ...
func (s *StockOpnameController) UploadTemplate(c *gin.Context) {
	// fmt.Println("-------->", req)

	var res dto.ContentResponse

	stockOpnameId, err := strconv.ParseInt(c.Param("stockOpnameId"), 10, 64)
	if err != nil {
		fmt.Println("Erro => ", err)
		res.ErrDesc = constants.ERR_CODE_05_MSG
		res.ErrCode = constants.ERR_CODE_05
		res.Data = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	file, header, err := c.Request.FormFile("data-stock-opname.csv")
	if err != nil {
		fmt.Println("Error", err)
		res.ErrDesc = constants.ERR_CODE_06_MSG
		res.ErrCode = constants.ERR_CODE_06
		res.Data = err.Error()
		c.JSON(http.StatusOK, res)

		return
	}
	filename := header.Filename
	fmt.Println("File name ==>", filename)
	templateRecords, errcode, errdesc, msg := convertCSVtoData(file)

	if errcode != constants.ERR_CODE_00 {
		res.ErrDesc = errdesc
		res.ErrCode = errcode
		res.Data = msg
		// log.Fatal(msg, err)
		c.JSON(http.StatusOK, res)
		return
	}

	// for _, rec := range templateRecords {
	// 	fmt.Println(rec.ProductName)
	// }

	fmt.Println("Set stock opname ", stockOpnameId)
	stockOpnameDetailService.SaveByUploadaData(stockOpnameId, templateRecords)

	res.ErrDesc = constants.ERR_CODE_00_MSG
	res.ErrCode = constants.ERR_CODE_00
	res.Data = ""
	c.JSON(http.StatusOK, res)
	return
}

func convertCSVtoData(file multipart.File) ([]dto.TemplateReportStockOpname, string, string, string) {

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {

		return nil, constants.ERR_CODE_07, constants.ERR_CODE_07_MSG, "Unable to parse file as CSV"
	}

	var datas []dto.TemplateReportStockOpname
	for _, value := range records {
		var template dto.TemplateReportStockOpname
		template.ProductID = util.Atoi64(value[0])
		template.ProductName = value[1]
		template.Qty = util.Atoi64(value[2])
		template.UomID = util.Atoi64(value[4])

		// uniq product ID
		addArray := true
		for _, dataCurr := range datas {
			if dataCurr.ProductID == template.ProductID {
				fmt.Println("already exist ", template.ProductID)
				addArray = false
				break
			}
		}
		if addArray {
			datas = append(datas, template)
		}
	}

	return datas, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, "ok"
}

func (s *StockOpnameController) RecalculateTotal(c *gin.Context) {

	res := dto.StockOpnameSaveResult{}
	errCode, errMsg := StockOpnameService.RecalculateTotal()
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)
	return
}
