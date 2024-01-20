package controllers

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
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

// PaymentController ...
type DirectPaymentController struct {
	DB *gorm.DB
}

// PaymentService ...
var DirectPaymentService = new(services.DirectPaymentService)

// FilterData ...
func (s *DirectPaymentController) FilterData(c *gin.Context) {
	req := dto.FilterPaymentDirect{}
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

	res = DirectPaymentService.GetDataPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *DirectPaymentController) Approve(c *gin.Context) {

	req := dto.PaymentDirectModel{}
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

	errCode, errMsg := DirectPaymentService.Approve(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *DirectPaymentController) Reject(c *gin.Context) {

	req := dto.PaymentDirectModel{}
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

	errCode, errMsg := DirectPaymentService.Reject(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}
