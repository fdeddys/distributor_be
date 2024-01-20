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

// PaymentOrderController ...
type PaymentOrderController struct {
	DB *gorm.DB
}

var paymentOrderService = new(services.PaymentOrderService)

// GetDetail ...
func (s *PaymentOrderController) GetDetail(c *gin.Context) {

	req := dto.FilterPayment{}
	res := models.ResponsePagination{}

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
	res = paymentOrderService.GetDataPaymentOrderPage(req)
	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *PaymentOrderController) Save(c *gin.Context) {

	req := dbmodels.PaymentOrder{}
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

	errCode, errMsg := paymentOrderService.SavePaymentOrder(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *PaymentOrderController) DeleteById(c *gin.Context) {

	res := models.NoContentResponse{}
	paymentOrderId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logs.Info("error", paymentOrderId)
		res.ErrDesc = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println(" id => ", paymentOrderId)

	res.ErrCode, res.ErrDesc = paymentOrderService.DeletePaymentOrderByID(paymentOrderId)
	c.JSON(http.StatusOK, res)

	return
}
