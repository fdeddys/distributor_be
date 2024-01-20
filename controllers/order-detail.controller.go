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

// OrderDetailController ...
type OrderDetailController struct {
	DB *gorm.DB
}

var orderDetailService = new(services.OrderDetailService)

// GetDetail ...
func (s *OrderDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterOrderDetail{}
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

	log.Println("order id => ", req.OrderID)

	res = orderDetailService.GetDataOrderDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *OrderDetailController) Save(c *gin.Context) {

	req := dbmodels.SalesOrderDetail{}
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

	errCode, errMsg := orderDetailService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *OrderDetailController) DeleteById(c *gin.Context) {

	res := dto.OrderDetailResult{}

	orderDetailId, errOrderId := strconv.ParseInt(c.Param("id"), 10, 64)
	if errOrderId != nil {
		logs.Info("error", errOrderId)
		res.ErrDesc = errOrderId.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("order id => ", orderDetailId)

	res.ErrCode, res.ErrDesc = orderDetailService.DeleteOrderDetailByID(orderDetailId)

	c.JSON(http.StatusOK, res)

	return
}

// UpdateQtyReceive ...
func (s *OrderDetailController) UpdateQtyReceive(c *gin.Context) {

	res := dto.OrderDetailResult{}
	req := dto.FilterOrderDetail{}
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

	res.ErrCode, res.ErrDesc = orderDetailService.UpdateQtyReceive(req.OrderDetailId, req.QtyReceive)

	c.JSON(http.StatusOK, res)

	return
}

// UpdateQtyOrder ...
func (s *OrderDetailController) UpdateQtyOrder(c *gin.Context) {

	res := dto.OrderDetailResult{}
	req := dto.FilterOrderDetail{}
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

	res.ErrCode, res.ErrDesc = orderDetailService.UpdateQtyOrder(req.OrderDetailId, req.QtyOrder)

	c.JSON(http.StatusOK, res)

	return
}
