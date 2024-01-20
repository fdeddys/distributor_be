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

// OrderController ...
type OrderController struct {
	DB *gorm.DB
}

// OrderService ...
var OrderService = new(services.OrderService)

// GetByOrderId ...
func (s *OrderController) GetByOrderId(c *gin.Context) {

	res := dbmodels.SalesOrder{}

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = OrderService.GetDataOrderById(orderID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// GetByOrderId ...
func (s *OrderController) CekTotal(c *gin.Context) {

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, 0)
		c.Abort()
		return
	}

	res := OrderService.GetTotalOrderById(orderID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// FilterData ...
func (s *OrderController) FilterData(c *gin.Context) {
	req := dto.FilterOrder{}
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
	res = OrderService.GetDataPage(req, page, count, intStatus)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *OrderController) Save(c *gin.Context) {

	req := dbmodels.SalesOrder{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, orderNo, orderID, status := OrderService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.OrderNo = orderNo
	res.ID = orderID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *OrderController) Approve(c *gin.Context) {

	req := dbmodels.SalesOrder{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := OrderService.Approve(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *OrderController) Reject(c *gin.Context) {

	req := dbmodels.SalesOrder{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := OrderService.Reject(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// PrintSO ...
func (s *OrderController) PrintSO(c *gin.Context) {

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GenerateSalesOrderReport(orderID, "so")

	header := c.Writer.Header()
	// header["Content-type"] = []string{"application/octet-stream"}
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename= tes.pdf"}

	// file, _ := os.Open("/Users/deddysyuhendra/go/src/tes-print/invoice.pdf")
	file, _ := os.Open("invoice.pdf")

	io.Copy(c.Writer, file)
	return
}

// PrintSO ...
func (s *OrderController) PrintInvoice(c *gin.Context) {

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GenerateSalesOrderReport(orderID, "invoice")

	header := c.Writer.Header()
	// header["Content-type"] = []string{"application/octet-stream"}
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename= tes.pdf"}

	// file, _ := os.Open("/Users/deddysyuhendra/go/src/tes-print/invoice.pdf")
	file, _ := os.Open("invoice.pdf")

	io.Copy(c.Writer, file)
	return
}

// CreateInvoice ...
func (s *OrderController) CreateInvoice(c *gin.Context) {

	req := dto.FilterOrder{}

	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Invoice Filter stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := OrderService.CreateInvoice(req.OrderID)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// FilterData ...
func (s *OrderController) FilterDataForSalesOrder(c *gin.Context) {
	req := dto.FilterOrder{}
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

	res = OrderService.GetDataForSalesOrderPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}
