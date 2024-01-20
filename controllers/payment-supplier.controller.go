package controllers

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PaymentSupplierController ...
type PaymentSupplierController struct {
	DB *gorm.DB
}

// PaymentService ...
var PaymentSupplierService = new(services.PaymentSupplierService)

// FilterData ...
func (s *PaymentSupplierController) FilterData(c *gin.Context) {
	req := dto.FilterSupplierPayment{}
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

	res = PaymentSupplierService.GetDataPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// GetByOrderId ...
func (s *PaymentSupplierController) GetByPaymentId(c *gin.Context) {

	res := dbmodels.PaymentSupplier{}

	paymentID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = PaymentSupplierService.GetDataById(paymentID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// Save ...
func (s *PaymentSupplierController) Save(c *gin.Context) {

	req := dbmodels.PaymentSupplier{}
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

	errCode, errMsg, paymentNo, paymentID, status := PaymentSupplierService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.PaymentNo = paymentNo
	res.ID = paymentID
	res.Status = status
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *PaymentSupplierController) Approve(c *gin.Context) {

	req := dbmodels.PaymentSupplier{}
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

	errCode, errMsg := PaymentSupplierService.Approve(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}
