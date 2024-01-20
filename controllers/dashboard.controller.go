package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"distribution-system-be/constants"
	"distribution-system-be/models"
	"distribution-system-be/models/dto"
	"distribution-system-be/services"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DashboardController ...
type DashboardController struct {
	DB *gorm.DB
}

// DashboardService ...
var DashboardService = new(services.DashboardService)

// FilterDataDashboard ...
func (i *DashboardController) FilterDataDashboard(c *gin.Context) {
	req := dto.FilterDto{}
	res := models.ContentResponse{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = constants.ERR_CODE_03
		res.ErrDesc = constants.ERR_CODE_03_MSG
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = DashboardService.GetQtyOrder(req)
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	c.JSON(http.StatusOK, res)
	return
}
