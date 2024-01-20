package controllers

import (
	"distribution-system-be/models"
	"distribution-system-be/services"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// StockController ...
type StockController struct {
	DB *gorm.DB
}

// StockMutationService ...
var StockService = new(services.StockService)

// GetByStockMutationId ...
func (s *StockController) GetByStockByProductIdPage(c *gin.Context) {

	res := models.ResponsePagination{}

	productId, errID := strconv.ParseInt(c.Param("id"), 10, 64)
	if errID != nil {
		logs.Info("error ID", errID)
		res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
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

	c.JSON(http.StatusOK, StockService.GetDataStockByProductPage(productId, page, count))
	c.Abort()
	return

}
