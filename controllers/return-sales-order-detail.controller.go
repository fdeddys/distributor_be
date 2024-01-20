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

// ReturnSalesOrderDetailController ...
type ReturnSalesOrderDetailController struct {
	DB *gorm.DB
}

var returnSalesOrderDetailService = new(services.ReturnOrderDetailService)

// GetDetail ...
func (s *ReturnSalesOrderDetailController) GetDetail(c *gin.Context) {

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

	log.Println("order id => ", req.OrderReturnID)

	res = returnSalesOrderDetailService.GetReturnOrderDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *ReturnSalesOrderDetailController) Save(c *gin.Context) {

	req := dbmodels.ReturnSalesOrderDetail{}
	body := c.Request.Body
	res := dto.OrderDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := returnSalesOrderDetailService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *ReturnSalesOrderDetailController) DeleteById(c *gin.Context) {

	res := dto.NoContentResponse{}

	orderDetailId, errOrderId := strconv.ParseInt(c.Param("id"), 10, 64)
	if errOrderId != nil {
		logs.Info("error", errOrderId)
		res.ErrDesc = errOrderId.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	res.ErrCode, res.ErrDesc = returnSalesOrderDetailService.DeleteReturnOrderDetailById(orderDetailId)

	c.JSON(http.StatusOK, res)

	return
}

// UpdateQtyReceive ...
func (s *ReturnSalesOrderDetailController) UpdateQty(c *gin.Context) {

	res := dto.OrderDetailResult{}
	req := dto.FilterOrderReturnDetail{}
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res.ErrCode, res.ErrDesc = returnSalesOrderDetailService.UpdateQtReturn(req.OrderReturnDetailId, req.QtyReturn)

	c.JSON(http.StatusOK, res)

	return
}
