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
type PaymentDetailController struct {
	DB *gorm.DB
}

var paymentDetailService = new(services.PaymentDetailService)

// GetDetail ...
func (s *PaymentDetailController) GetDetail(c *gin.Context) {

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
	res = paymentDetailService.GetDataPaymentDetailPage(req)
	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *PaymentDetailController) Save(c *gin.Context) {

	req := dbmodels.PaymentDetail{}
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

	errCode, errMsg := paymentDetailService.SavePaymentDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *PaymentDetailController) DeleteById(c *gin.Context) {

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

	res.ErrCode, res.ErrDesc = paymentDetailService.DeletePaymentDetailByID(paymentDetailId)
	c.JSON(http.StatusOK, res)

	return
}
