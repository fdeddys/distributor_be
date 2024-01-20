package controllers

import (
	"distribution-system-be/models/dto"
	"distribution-system-be/services/reportservice"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// StockOpnameController ...
type ReportPaymentSupplierController struct {
	DB *gorm.DB
}

// StockOpnameService ...
var reporPaymentSupplier = new(reportservice.ReportPaymentSupplierService)

// DownloadTemplate ...
func (s *ReportPaymentSupplierController) DownloadReporPaymentSupplier(c *gin.Context) {
	// fmt.Println("-------->", req)
	req := dto.FilterReportDate{}
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, unmarshal body Request to Receive stuct ", dataBodyReq)
		c.JSON(http.StatusBadRequest, "")
		c.Abort()
		return
	}

	filename, success := reporPaymentSupplier.GenerateReportPaymentSupplier(req)
	if success {
		fmt.Println("download PDF ")
		header := c.Writer.Header()
		header["Content-type"] = []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}
		header["Content-Disposition"] = []string{"attachment; filename=" + filename}

		file, _ := os.Open(filename)

		io.Copy(c.Writer, file)
		os.Remove(filename)
	}
	c.JSON(http.StatusOK, "Success !")

	return
}
