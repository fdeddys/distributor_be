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

// PurchaseOrderDetailController ...
type PurchaseOrderDetailController struct {
	DB *gorm.DB
}

var purchaseOrderDetailService = new(services.PurchaseOrderDetailService)

// GetDetail ...
func (r *PurchaseOrderDetailController) GetDetail(c *gin.Context) {

	req := dto.FilterPurchaseOrderDetail{}
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

	log.Println("order id => ", req.PurchaseOrderID)

	res = purchaseOrderDetailService.GetDataPurchaseOrderDetailPage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (r *PurchaseOrderDetailController) Save(c *gin.Context) {

	req := dbmodels.PurchaseOrderDetail{}
	body := c.Request.Body
	res := dto.PurchaseOrderDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", err)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := purchaseOrderDetailService.SavePurchaseOrderDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// DeleteByID ...
func (r *PurchaseOrderDetailController) DeleteByID(c *gin.Context) {

	res := dto.NoContentResponse{}

	purchaseOrderDetailID, errPurchaseOrderID := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPurchaseOrderID != nil {
		logs.Info("error", purchaseOrderDetailID)
		res.ErrDesc = errPurchaseOrderID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("purchase Order id => ", purchaseOrderDetailID)
	res.ErrCode, res.ErrDesc = purchaseOrderDetailService.DeletePurchaseOrderDetailByID(purchaseOrderDetailID)
	c.JSON(http.StatusOK, res)

	return
}

// GetDetail ...
func (r *PurchaseOrderDetailController) GetLastPrice(c *gin.Context) {

	res := models.ResponseCheckPrice{}

	productId, err := strconv.ParseInt(c.Param("productId"), 10, 64)

	if err != nil {
		logs.Info("error", err)
		res.ErrCode = constants.ERR_CODE_03
		res.ErrDesc = constants.ERR_CODE_03_MSG
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("productCode->", productId)
	res = purchaseOrderDetailService.GetLastPricePurchaseOrderDetail(productId)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (r *PurchaseOrderDetailController) UpdateDetail(c *gin.Context) {

	req := dbmodels.PurchaseOrderDetail{}
	body := c.Request.Body
	res := dto.PurchaseOrderDetailSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", err)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := purchaseOrderDetailService.UpdateDetail(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}
