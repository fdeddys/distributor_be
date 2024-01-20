package controllers

import (
	"distribution-system-be/services/reportservice"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ReportMasterController ...
type ReportMasterController struct {
	DB *gorm.DB
}

// ReportMasterBarangService ...
var reportMasterBarangService = new(reportservice.ReportMasterService)

// DownloadTemplate ...
func (s *ReportMasterController) DownloadReportMasterProduct(c *gin.Context) {

	filename, success := reportMasterBarangService.GenerateReportMasterProduct()
	if success {
		header := c.Writer.Header()
		header["Content-type"] = []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}
		header["Content-Disposition"] = []string{"attachment; filename=" + filename}

		file, _ := os.Open(filename)

		io.Copy(c.Writer, file)
		os.Remove(filename)
	}
	c.JSON(http.StatusOK, "Success !")

	// header := c.Writer.Header()
	// header["Content-type"] = []string{"text/csv"}
	// header["Content-Disposition"] = []string{"attachment; filename=report.csv"}

	// file, _ := os.Open(filename)

	// io.Copy(c.Writer, file)

	// os.Remove(filename)
	return
}
