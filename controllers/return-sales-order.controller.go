package controllers

import (
	"distribution-system-be/constants"
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

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RetusnSalesOrderController ...
type RetusnSalesOrderController struct {
	DB *gorm.DB
}

// ReturnSalesOrderService ...
var ReturnSalesOrderService = new(services.ReturnSalesOrderService)

// GetByOrderId ...
func (s *RetusnSalesOrderController) GetByReturnSalesOrderId(c *gin.Context) {

	res := dbmodels.ReturnSalesOrder{}

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ReturnSalesOrderService.GetDataOrderReturnById(orderID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// FilterData ...
func (s *RetusnSalesOrderController) FilterData(c *gin.Context) {
	req := dto.FilterOrderReturnDetail{}
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

	res = ReturnSalesOrderService.GetDataSalesOrderReturnPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *RetusnSalesOrderController) Save(c *gin.Context) {

	req := dbmodels.ReturnSalesOrder{}
	body := c.Request.Body
	res := dto.ReturnOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, orderNo, returnOrderID, status := ReturnSalesOrderService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.ReturnNo = orderNo
	res.ID = returnOrderID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *RetusnSalesOrderController) Approve(c *gin.Context) {

	req := dbmodels.ReturnSalesOrder{}
	body := c.Request.Body
	res := dto.ReturnOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := ReturnSalesOrderService.Approve(req.ID)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *RetusnSalesOrderController) Reject(c *gin.Context) {

	req := dbmodels.ReturnSalesOrder{}
	body := c.Request.Body
	res := dto.ReturnOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := ReturnSalesOrderService.Reject(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// PrintSO ...
func (s *RetusnSalesOrderController) PrintReturnSO(c *gin.Context) {

	returnSoID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GenerateReturnSalesOrderReport(returnSoID)

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename=tes.pdf"}

	file, _ := os.Open("return-so.pdf")

	io.Copy(c.Writer, file)
	return
}

// FilterData ...
func (s *RetusnSalesOrderController) FilterDataForSalesOrderReturn(c *gin.Context) {
	req := dto.FilterOrderReturnDetail{}
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

	// temp, _ := json.Marshal(req)
	// log.Println("searchName-->", string(temp))
	// log.Println("is release ", req.InternalStatus)

	res = ReturnSalesOrderService.GetDataForSalesOrderReturnPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}
