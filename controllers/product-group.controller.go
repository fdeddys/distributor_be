package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/services"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ProductGroupController ...
type ProductGroupController struct {
	DB *gorm.DB
}

//ProductGroupService ...
var ProductGroupService = new(services.ProductGroupService)

//GetProductGroupPaging ...
func (h *ProductGroupController) GetProductGroupPaging(c *gin.Context) {
	req := dto.FilterProductGroup{}
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
		fmt.Println("Error, body Request ")
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ProductGroupService.GetProductGroupPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetProductGroupDetails ...
func (h *ProductGroupController) GetProductGroupDetails(c *gin.Context) {
	res := models.ContentResponse{}

	id, errID := strconv.Atoi(c.Param("id"))
	if errID != nil {
		logs.Info("error", errID)
		// res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ProductGroupService.GetProductGroupDetails(id)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// SaveProductGroup ...
func (h *ProductGroupController) SaveProductGroup(c *gin.Context) {

	req := dbmodels.ProductGroup{}
	res := models.NoContentResponse{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = "03"
		res.ErrDesc = "Error, unmarshall body Request"
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, ProductGroupService.SaveProductGroup(&req))
	return
}

// UpdateProductGroup ...
func (h *ProductGroupController) UpdateProductGroup(c *gin.Context) {
	req := dbmodels.ProductGroup{}
	res := models.NoContentResponse{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = "03"
		res.ErrDesc = "Error, unmarshall body Request"
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, ProductGroupService.UpdateProductGroup(&req))
	return
}
