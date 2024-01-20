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

// ReturnReceiveController ...
type ReturnReceiveController struct {
	DB *gorm.DB
}

// ReturnReceiveService ...
var ReturnReceiveService = new(services.ReturnReceiveService)

// GetByOrderId ...
func (s *ReturnReceiveController) GetByReturnReceiveId(c *gin.Context) {

	res := dbmodels.ReturnReceive{}

	returnReceiveID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ReturnReceiveService.GetDataReturnById(returnReceiveID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// FilterData ...
func (s *ReturnReceiveController) FilterData(c *gin.Context) {
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

	temp, _ := json.Marshal(req)
	log.Println("searchName-->", string(temp))

	res = ReturnReceiveService.GetDataReturnReceivePage(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (s *ReturnReceiveController) Save(c *gin.Context) {

	req := dbmodels.ReturnReceive{}
	body := c.Request.Body
	res := dto.ReturnOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Return Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, orderNo, returnOrderID, status := ReturnReceiveService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.ReturnNo = orderNo
	res.ID = returnOrderID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *ReturnReceiveController) Approve(c *gin.Context) {

	req := dbmodels.ReturnReceive{}
	body := c.Request.Body
	res := dto.ReturnOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Return Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := ReturnReceiveService.Approve(req.ID)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *ReturnReceiveController) Reject(c *gin.Context) {

	req := dbmodels.ReturnReceive{}
	body := c.Request.Body
	res := dto.ReturnOrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Return Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := ReturnReceiveService.Reject(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

// PrintSO ...
func (s *ReturnReceiveController) PrintReturn(c *gin.Context) {

	returnReceiveID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GenerateReturnReceiveReport(returnReceiveID)

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename=tes.pdf"}

	file, _ := os.Open("return-receive.pdf")

	io.Copy(c.Writer, file)
	return
}
