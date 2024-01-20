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

// ReceiveDetailController ...
type ReceiveDetailController struct {
	DB *gorm.DB
}

var receiveDetailService = new(services.ReceiveDetailService)

// GetDetail ...
func (r *ReceiveDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterReceiveDetail{}
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

	log.Println("order id => ", req.ReceiveID)

	res = receiveDetailService.GetDataReceiveDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (r *ReceiveDetailController) Save(c *gin.Context) {

	req := dbmodels.ReceiveDetail{}
	body := c.Request.Body
	res := dto.ReceiveDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := receiveDetailService.SaveReceiveDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (r *ReceiveDetailController) UpdateDetail(c *gin.Context) {

	req := []dbmodels.ReceiveDetail{}
	body := c.Request.Body
	res := dto.ReceiveDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := receiveDetailService.UpdateReceiveDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteByID Multiple ...
func (r *ReceiveDetailController) DeleteByIDMultiple(c *gin.Context) {

	req := []int64{}
	body := c.Request.Body
	res := dto.ReceiveDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	 err := json.Unmarshal(dataBodyReq, &req)
	 if err != nil {
		fmt.Println("Error, unmarshal body Request  ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	fmt.Println("ISI " , req)

	res.ErrCode, res.ErrDesc = receiveDetailService.DeleteReceiveDetailByIDMultiple(req)

	c.JSON(http.StatusOK, res)

	return
}


// DeleteByID ...
func (r *ReceiveDetailController) DeleteByID(c *gin.Context) {

	res := dto.ReceiveDetailResult{}

	receiveDetailID, errReceiveID := strconv.ParseInt(c.Param("id"), 10, 64)
	if errReceiveID != nil {
		logs.Info("error", receiveDetailID)
		res.ErrDesc = errReceiveID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("receive id => ", receiveDetailID)

	res.ErrCode, res.ErrDesc = receiveDetailService.DeleteReceiveDetailByID(receiveDetailID)

	c.JSON(http.StatusOK, res)

	return
}

// GetDetail ...
func (r *ReceiveDetailController) SearchBatchExpired(c *gin.Context) {

	req := dto.FilterBatchExpired{}
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
		fmt.Println("Error, body Request ", res)
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = receiveDetailService.GetDataBatchExpired(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// GetDetail ...
func (r *ReceiveDetailController) GetLastPrice(c *gin.Context) {

	res := models.ResponseReceiveCheckPrice{}
	productId, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		logs.Info("error", err)
		res.Price = 0
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	res = receiveDetailService.GetDataPriceProduct(productId)

	c.JSON(http.StatusOK, res)

	return
}
