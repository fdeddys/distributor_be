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

// ReturnReceiveDetailController ...
type ReturnReceiveDetailController struct {
	DB *gorm.DB
}

var returnReceiveDetailService = new(services.ReturnReceiveDetailService)

// GetDetail ...
func (s *ReturnReceiveDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterReturnReceive{}
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

	log.Println("return id => ", req.ReturnID)

	res = returnReceiveDetailService.GetReturnReceiveDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *ReturnReceiveDetailController) Save(c *gin.Context) {

	req := dbmodels.ReturnReceiveDetail{}
	body := c.Request.Body
	res := dto.OrderDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Return receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := returnReceiveDetailService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *ReturnReceiveDetailController) DeleteById(c *gin.Context) {

	res := dto.NoContentResponse{}

	detailId, errReceive := strconv.ParseInt(c.Param("id"), 10, 64)
	if errReceive != nil {
		logs.Info("error", errReceive)
		res.ErrDesc = errReceive.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	res.ErrCode, res.ErrDesc = returnReceiveDetailService.DeleteReturnReceiveDetailById(detailId)

	c.JSON(http.StatusOK, res)

	return
}

// UpdateQtyReceive ...
func (s *ReturnReceiveDetailController) UpdateQty(c *gin.Context) {

	res := dto.OrderDetailResult{}
	req := dto.FilterReturnReceive{}
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Return receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res.ErrCode, res.ErrDesc = returnSalesOrderDetailService.UpdateQtReturn(req.ReturnDetailId, req.Qty)

	c.JSON(http.StatusOK, res)

	return
}
