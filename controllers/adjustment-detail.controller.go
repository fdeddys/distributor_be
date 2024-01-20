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

// AdjustmentDetailController ...
type AdjustmentDetailController struct {
	DB *gorm.DB
}

var adjustmentDetailService = new(services.AdjustmentDetailService)

// GetDetail ...
func (r *AdjustmentDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterAdjustmentDetail{}
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

	log.Println("order id => ", req.AdjustmentID)

	res = adjustmentDetailService.GetDataAdjustmentDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (r *AdjustmentDetailController) Save(c *gin.Context) {

	req := dbmodels.AdjustmentDetail{}
	body := c.Request.Body
	res := dto.AdjustmentDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := adjustmentDetailService.SaveAdjustmentDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteByID ...
func (r *AdjustmentDetailController) DeleteByID(c *gin.Context) {

	res := dto.AdjustmentDetailResult{}

	adjustmentDetailID, errAdjustmentID := strconv.ParseInt(c.Param("id"), 10, 64)
	if errAdjustmentID != nil {
		logs.Info("error", adjustmentDetailID)
		res.ErrDesc = errAdjustmentID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	adjustmentID, errAdjustmentID := strconv.ParseInt(c.Param("idAdj"), 10, 64)

	log.Println("adjustment id => ", adjustmentDetailID)

	res.ErrCode, res.ErrDesc = adjustmentDetailService.DeleteAdjustmentDetailByID(adjustmentDetailID,adjustmentID)

	c.JSON(http.StatusOK, res)

	return
}



// UpdatQty ...
func (r *AdjustmentDetailController) UpdatQty(c *gin.Context) {

	res := dto.AdjustmentDetailResult{}

	req := dbmodels.AdjustmentDetail{}
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Adj stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res.ErrCode, res.ErrDesc = adjustmentDetailService.UpdateQtyByID(req.ID, req.Qty, req.AdjustmentID)

	c.JSON(http.StatusOK, res)

	return
}
