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

// AdjustmentController ...
type AdjustmentController struct {
	DB *gorm.DB
}

// AdjustmentService ...
var AdjustmentService = new(services.AdjustmentService)

// GetByAdjustmentId ...
func (a *AdjustmentController) GetByAdjustmentId(c *gin.Context) {

	res := dbmodels.Adjustment{}

	adjustmentID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = AdjustmentService.GetDataAdjustmentByID(adjustmentID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// FilterData ...
func (a *AdjustmentController) FilterData(c *gin.Context) {
	req := dto.FilterAdjustment{}
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
	log.Println("is release ", req.Status)

	intStatus := -1
	if intVal, errconv := strconv.Atoi(req.Status); errconv == nil {
		intStatus = intVal
	}
	res = AdjustmentService.GetDataPage(req, page, count, intStatus)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (a *AdjustmentController) Save(c *gin.Context) {

	req := dbmodels.Adjustment{}
	body := c.Request.Body
	res := dto.AdjustmentSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Adjustment stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, adjustmentNo, adjustmentID, status := AdjustmentService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.AdjustmentNo = adjustmentNo
	res.ID = adjustmentID
	res.Status = status
	// res.AdjustmentNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (a *AdjustmentController) Approve(c *gin.Context) {

	req := dbmodels.Adjustment{}
	body := c.Request.Body
	res := dto.AdjustmentSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Adjustment stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := AdjustmentService.ApproveAdjustment(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.AdjustmentNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (a *AdjustmentController) Reject(c *gin.Context) {

	req := dbmodels.Adjustment{}
	body := c.Request.Body
	res := dto.AdjustmentSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Adjustment stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := AdjustmentService.RejectAdjustment(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.AdjustmentNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// PrintPreview ...
func (a *AdjustmentController) PrintPreview(c *gin.Context) {

	adjustmentID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	fmt.Println("-------->", adjustmentID)

	report.GenerateSalesAdjustmentReport(adjustmentID)

	header := c.Writer.Header()
	// header["Content-type"] = []string{"application/octet-stream"}
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename= tes.pdf"}

	// file, _ := os.Open("/Users/deddysyuhendra/go/src/tes-print/invoice.pdf")
	file, _ := os.Open("adjustment.pdf")

	io.Copy(c.Writer, file)
	return
}
