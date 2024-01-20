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

// PaymentDetailController ...
type PaymentSupplierDetailController struct {
	DB *gorm.DB
}

var paymentSupplierDetailService = new(services.PaymentSupplierDetailService)

// GetDetail ...
func (s *PaymentSupplierDetailController) GetDetail(c *gin.Context) {

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
		// res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("id => ", req.PaymentID)
	res = paymentSupplierDetailService.GetDataPaymentDetailPage(req, page, count)
	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *PaymentSupplierDetailController) Save(c *gin.Context) {

	req := dbmodels.PaymentSupplierDetail{}
	body := c.Request.Body
	res := models.ContentResponse{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := paymentSupplierDetailService.SavePaymentDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *PaymentSupplierDetailController) DeleteById(c *gin.Context) {

	res := models.NoContentResponse{}
	paymentDetailId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.Info("error", paymentDetailId)
		res.ErrDesc = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println(" id => ", paymentDetailId)

	res.ErrCode, res.ErrDesc = paymentSupplierDetailService.DeletePaymentDetailByID(paymentDetailId)
	c.JSON(http.StatusOK, res)

	return
}
