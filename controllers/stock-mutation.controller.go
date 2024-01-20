package controllers

import (
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/report"
	"distribution-system-be/services"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"distribution-system-be/constants"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// StockMutationController ...
type StockMutationController struct {
	DB *gorm.DB
}

// StockMutationService ...
var StockMutationService = new(services.StockMutationService)

// GetByStockMutationId ...
func (s *StockMutationController) GetByStockMutationById(c *gin.Context) {

	res := dbmodels.StockMutation{}

	stockMutationID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = StockMutationService.GetDataStockMutationById(stockMutationID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// GetStockMutationPage ...
func (s *StockMutationController) GetStockMutationPage(c *gin.Context) {
	req := dto.FilterStockMutation{}
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
	res = StockMutationService.GetDataPage(req, page, count, intStatus)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *StockMutationController) Save(c *gin.Context) {

	req := dbmodels.StockMutation{}
	body := c.Request.Body
	res := dto.StockMutationSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales StockMutation stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, stockMutationNo, stockMutationID, status := StockMutationService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.StockMutationNo = stockMutationNo
	res.ID = stockMutationID
	res.Status = status
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *StockMutationController) Approve(c *gin.Context) {

	req := dbmodels.StockMutation{}
	body := c.Request.Body
	res := dto.StockMutationSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales StockMutation stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		return
	}

	errCode, errMsg := StockMutationService.Approve(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)
	return
}

// Approve ...
func (s *StockMutationController) Reject(c *gin.Context) {

	req := dbmodels.StockMutation{}
	body := c.Request.Body
	res := dto.StockMutationSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales StockMutation stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		return
	}

	errCode, errMsg := StockMutationService.Reject(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)
	return
}

// PrintSO ...
func (s *StockMutationController) PrintStockMutationForm(c *gin.Context) {

	stockMutationID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GenerateStockMutationReport(stockMutationID)

	header := c.Writer.Header()
	// header["Content-type"] = []string{"application/octet-stream"}
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename= tes.pdf"}

	// file, _ := os.Open("/Users/deddysyuhendra/go/src/tes-print/invoice.pdf")
	file, _ := os.Open("stock-mutation.pdf")

	io.Copy(c.Writer, file)
	return
}
