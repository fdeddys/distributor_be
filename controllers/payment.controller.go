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

// PaymentController ...
type PaymentController struct {
	DB *gorm.DB
}

// PaymentService ...
var PaymentService = new(services.PaymentService)

// GetByOrderId ...
func (s *PaymentController) GetByPaymentId(c *gin.Context) {

	res := dbmodels.Payment{}

	paymentID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = PaymentService.GetDataById(paymentID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// GetByOrderId ...
func (s *PaymentController) GetBySalesOrderID(c *gin.Context) {

	res := dbmodels.Payment{}

	salesOrderID, errPage := strconv.ParseInt(c.Param("salesOrderId"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = PaymentService.GetDataPaymentBySalesOrderId(salesOrderID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// FilterData ...
func (s *PaymentController) FilterData(c *gin.Context) {
	req := dto.FilterPayment{}
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

	res = PaymentService.GetDataPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *PaymentController) Save(c *gin.Context) {

	req := dbmodels.Payment{}
	body := c.Request.Body
	res := dto.SavePaymentResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to PAYMENT stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, paymentNo, paymentID, status := PaymentService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.PaymentNo = paymentNo
	res.ID = paymentID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *PaymentController) Approve(c *gin.Context) {

	req := dbmodels.Payment{}
	body := c.Request.Body
	res := dto.SaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Payment stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := PaymentService.Approve(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *PaymentController) Reject(c *gin.Context) {

	req := dbmodels.Payment{}
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

	errCode, errMsg := PaymentService.Reject(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// PrintSO ...
func (s *PaymentController) PrintPayment(c *gin.Context) {

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
