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

// StockOpnameDetailController ...
type StockOpnameDetailController struct {
	DB *gorm.DB
}

var stockOpnameDetailService = new(services.StockOpnameDetailService)

// GetDetail ...
func (s *StockOpnameDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterStockOpname{}
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

	log.Println("order id => ", req.StockOpnameID)

	res = stockOpnameDetailService.GetDataStockOpnameDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *StockOpnameDetailController) Save(c *gin.Context) {

	req := dbmodels.StockOpnameDetail{}
	body := c.Request.Body
	res := dto.StockOpnameSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := stockOpnameDetailService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *StockOpnameDetailController) DeleteById(c *gin.Context) {

	res := dto.NoContentResponse{}

	detailId, errOrderId := strconv.ParseInt(c.Param("id"), 10, 64)
	if errOrderId != nil {
		logs.Info("error", errOrderId)
		res.ErrDesc = errOrderId.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("detail Id  => ", detailId)

	res.ErrCode, res.ErrDesc = stockOpnameDetailService.DeleteStockOpnameDetailByID(detailId)

	c.JSON(http.StatusOK, res)

	return
}

// UpdateQty ...
func (s *StockOpnameDetailController) UpdateQty(c *gin.Context) {

	res := dto.StockOpnameSaveResult{}
	req := dto.FilterStockOpname{}
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

	res.ErrCode, res.ErrDesc = stockOpnameDetailService.UpdateQty(req.StockOpnameDetailId, req.Qty)

	c.JSON(http.StatusOK, res)

	return
}
