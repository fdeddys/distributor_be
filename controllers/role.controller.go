package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	constants "distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/services"
	"strconv"

	dto "distribution-system-be/models/dto"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RoleController ...
type RoleController struct {
	DB *gorm.DB
}

// RoleService ...
var RoleService = new(services.RoleService)

// GetRole ...
func (h *RoleController) GetRole(c *gin.Context) {
	req := dto.FilterName{}
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
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = RoleService.GetRoleFilterPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// SaveRole ...
func (h *RoleController) SaveRole(c *gin.Context) {

	req := dbmodels.Role{}
	res := models.NoContentResponse{}

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

	c.JSON(http.StatusOK, RoleService.SaveRole(&req))
	return
}

// UpdateRole ...
func (h *RoleController) UpdateRole(c *gin.Context) {
	req := dbmodels.Role{}
	res := models.NoContentResponse{}

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

	c.JSON(http.StatusOK, RoleService.UpdateRole(&req))
	return
}
