package controllers

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/report"
	"distribution-system-be/services"
	"io"
	"os"

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

// ReceiveController ...
type ReceiveController struct {
	DB *gorm.DB
}

var receiveService = new(services.ReceiveService)

// FilterData ...
func (r *ReceiveController) FilterData(c *gin.Context) {
	req := dto.FilterReceive{}
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

	// status := -1
	// if intVal, errconv := strconv.Atoi(req.Status); errconv == nil {
	// 	status = intVal
	// }
	status := req.Status

	res = receiveService.GetDataPage(req, page, count, status)

	fmt.Println("res page =>", res)
	c.JSON(http.StatusOK, res)

	return
}

// GetByReceiveId ...
func (r *ReceiveController) GetByReceiveId(c *gin.Context) {

	res := dbmodels.Receive{}

	orderID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = receiveService.GetDataReceiveByID(orderID)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// Save ...
func (r *ReceiveController) Save(c *gin.Context) {

	req := dbmodels.Receive{}
	body := c.Request.Body
	res := dto.ReceiveSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg, receiveNo, receiveID, status := receiveService.Save(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.ReceiveNo = receiveNo
	res.ID = receiveID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Save ...
func (r *ReceiveController) SaveByPO(c *gin.Context) {

	req := dbmodels.Receive{}
	body := c.Request.Body
	res := dto.ReceiveSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Receive stuct ", err.Error())
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	fmt.Println("Req ==>", req)
	errCode, errMsg, receiveNo, receiveID, status := receiveService.SaveByPO(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	res.ReceiveNo = receiveNo
	res.ID = receiveID
	res.Status = status
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (r *ReceiveController) Approve(c *gin.Context) {

	req := dbmodels.Receive{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := receiveService.ApproveReceive(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// PrintPreview ...
func (r *ReceiveController) PrintPreview(c *gin.Context) {

	receiveID, errPage := strconv.ParseInt(c.Param("id"), 10, 64)
	if errPage != nil {
		logs.Info("error", errPage)
		c.JSON(http.StatusBadRequest, "id not supplied")
		c.Abort()
		return
	}

	// fmt.Println("-------->", req)

	report.GenerateReceiveReport(receiveID)

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/x-pdf"}
	header["Content-Disposition"] = []string{"attachment; filename= tes.pdf"}

	file, _ := os.Open("receive.pdf")

	io.Copy(c.Writer, file)
	return
}

// RemovePO ...
func (r *ReceiveController) RemovePO(c *gin.Context) {

	req := dbmodels.Receive{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := receiveService.RemovePO(&req, false)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

func (r *ReceiveController) RemovePOAllItem(c *gin.Context) {

	req := dbmodels.Receive{}
	body := c.Request.Body
	res := dto.OrderSaveResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Receive stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := receiveService.RemovePO(&req, true)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	// res.OrderNo = newNumb
	c.JSON(http.StatusOK, res)

	return
}

// Approve ...
func (s *ReceiveController) Reject(c *gin.Context) {

	req := dbmodels.Receive{}
	body := c.Request.Body
	res := dto.ReceiveDetailResult{}
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Sales Order stuct ", dataBodyReq)
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.ErrCode = constants.ERR_CODE_03
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	errCode, errMsg := receiveService.RejectReceive(&req)
	res.ErrDesc = errMsg
	res.ErrCode = errCode
	c.JSON(http.StatusOK, res)

	return
}

func (r *ReceiveController) Export(c *gin.Context) {

	req := dto.FilterReceive{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ", err)
		c.JSON(http.StatusBadRequest, "Failed un marshal")
		c.Abort()
		return
	}

	// temp, _ := json.Marshal(req)

	success, filename := receiveService.ExportReceive(req, req.Status)

	if success {
		header := c.Writer.Header()
		header["Content-type"] = []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}
		header["Content-Disposition"] = []string{"attachment; filename=" + filename}
		file, _ := os.Open(filename)

		io.Copy(c.Writer, file)

	}
	c.JSON(http.StatusOK, "Failed !")

	return
}
