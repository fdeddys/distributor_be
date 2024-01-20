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

// StockMutationDetailController ...
type StockMutationDetailController struct {
	DB *gorm.DB
}

var stockMutationDetailService = new(services.StockMutationDetailService)

// GetDetail ...
func (s *StockMutationDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterStockMutation{}
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

	log.Println("order id => ", req.StockMutationID)

	res = stockMutationDetailService.GetDataStockMutationDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *StockMutationDetailController) Save(c *gin.Context) {

	req := dbmodels.StockMutationDetail{}
	body := c.Request.Body
	res := dto.StockMutationSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := stockMutationDetailService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// DeleteById ...
func (s *StockMutationDetailController) DeleteById(c *gin.Context) {

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

	res.ErrCode, res.ErrDesc = stockMutationDetailService.DeleteStockMutationDetailByID(detailId)

	c.JSON(http.StatusOK, res)

	return
}

// UpdateQty ...
func (s *StockMutationDetailController) UpdateQty(c *gin.Context) {

	res := dto.StockMutationSaveResult{}
	req := dto.FilterStockMutation{}
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

	res.ErrCode, res.ErrDesc = stockMutationDetailService.UpdateQty(req.StockMutationDetailId, req.Qty)

	c.JSON(http.StatusOK, res)

	return
}
