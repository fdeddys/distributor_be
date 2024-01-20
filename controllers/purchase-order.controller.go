package controllers

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/report"
	"distribution-system-be/services"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PurchaseOrderController ...
type PurchaseOrderController struct {
	DB *gorm.DB
}

var purchaseOrderService = new(services.PurchaseOrderService)

// FilterData ...
func (r *PurchaseOrderController) FilterData(c *gin.Context) {
	req := dto.FilterPurchaseOrder{}
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
		fmt.Println("Error, body Request ", err)
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	temp, _ := json.Marshal(req)
	log.Println("searchName-->", string(temp))

	// status := -1
	// if intVal, errconv := strconv.Atoi(req.Status); errconv == nil {
	// 	status = intVal
	// }
	status := req.Status
	res = purchaseOrderService.GetDataPage(req, page, count, status)

	c.JSON(http.StatusOK, res)

	return
}

// GetByPurchaseOrderId ...
func (r *PurchaseOrderController) GetByPurchaseOrderId(c *gin.Context) {

	res := dbmodels.PurchaseOrder{}

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = purchaseOrderService.GetDataPurchaseOrderByID(orderID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// Save ...
func (r *PurchaseOrderController) Save(c *gin.Context) {

	req := dbmodels.PurchaseOrder{}
	body := c.Request.Body
	res := dto.PurchaseOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to PurchaseOrder stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, purchaseOrderNo, purchaseOrderID, status := purchaseOrderService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.PurchaseOrderNo = purchaseOrderNo
	res.ID = purchaseOrderID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (r *PurchaseOrderController) Approve(c *gin.Context) {

	req := dbmodels.PurchaseOrder{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to PurchaseOrder stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := purchaseOrderService.ApprovePurchaseOrder(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// PrintPreview ...
func (r *PurchaseOrderController) PrintPreview(c *gin.Context) {

	purchaseOrderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GeneratePurchaseOrderReport(purchaseOrderID)

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename=po.pdf"}

	file, _ := os.Open("purchase-order.pdf")

	io.Copy(c.Writer, file)
	return
}

// PrintPreviewByPoNo ...
func (r *PurchaseOrderController) PrintPreviewByPoNo(c *gin.Context) {

	pono := c.Param("pono")

	err := report.GeneratePurchaseOrderReportByPoNo(pono)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}
	header := c.Writer.Header()
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename=po.pdf"}

	file, _ := os.Open("purchase-order.pdf")

	io.Copy(c.Writer, file)
	return
}

// Approve ...
func (r *PurchaseOrderController) Reject(c *gin.Context) {

	res := dbmodels.PurchaseOrder{}

	poID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, purchaseOrderService.RejectPO(poID))

	return
}

// Approve ...
func (r *PurchaseOrderController) CancelSubmit(c *gin.Context) {

	res := dbmodels.PurchaseOrder{}

	poID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, purchaseOrderService.CancelSubmitPO(poID))

	return
}

func (r *PurchaseOrderController) Export(c *gin.Context) {

	req := dto.FilterPurchaseOrder{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ", err)
		c.JSON(http.StatusBadRequest, "Failed un marshal")
		c.Abort()
		return
	}

	// temp, _ := json.Marshal(req)

	success, filename := purchaseOrderService.ExportPurchaseOrder(req, req.Status)

	if success {
		fmt.Println("download template")
		header := c.Writer.Header()
		header["Content-type"] = []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}
		header["Content-Disposition"] = []string{"attachment; filename=" + filename}

		file, _ := os.Open(filename)

		io.Copy(c.Writer, file)

	}
	c.JSON(http.StatusOK, "Failed !")

	return
}
