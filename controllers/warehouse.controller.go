package controllers

import (
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// WarehouseController ...
type WarehouseController struct {
	DB *gorm.DB
}

// WarehouseService ...
var WarehouseService = new(services.WarehouseService)

// GetWarehouse ...
func (h *WarehouseController) GetWarehouse(c *gin.Context) {

	res := WarehouseService.GetAllWarehouse()

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetWarehouse IN ...
func (h *WarehouseController) GetWarehouseIn(c *gin.Context) {

	res := WarehouseService.GetWarehouseByFunc(true)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetWarehouse IN ...
func (h *WarehouseController) GetWarehouseOut(c *gin.Context) {

	res := WarehouseService.GetWarehouseByFunc(false)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetWarehouseFilter ...
func (h *WarehouseController) GetWarehouseFilter(c *gin.Context) {
	req := dto.FilterPaging{}
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

	res = WarehouseService.GetWarehouseFilterPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetFilterBrand ...
func (h *WarehouseController) GetFilterWarehouse(c *gin.Context) {
	res := models.ContentResponse{}

	id, errID := strconv.Atoi(c.Param("id"))
	if errID != nil {
		logs.Info("error", errID)
		// res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = WarehouseService.GetWarehouseFilter(id)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

//GetWarehouseLike ...
func (h *WarehouseController) GetWarehouseLike(c *gin.Context) {
	res := models.ContentResponse{}

	terms := c.Query("terms")

	if terms == "" {
		logs.Info("error", "can't found the warehouse string")
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	// fmt.Sprintf("ini lookupstr = " + lookupstr)

	res = WarehouseService.GetWarehouseLike(terms)
	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// Save Salaeman ...
func (h *WarehouseController) SaveWarehouse(c *gin.Context) {

	req := dbmodels.Warehouse{}
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

	c.JSON(http.StatusOK, WarehouseService.SaveWarehouse(&req))
	return
}

// UpdateBrand ...
func (h *WarehouseController) UpdateWarehouse(c *gin.Context) {
	req := dbmodels.Warehouse{}
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

	c.JSON(http.StatusOK, WarehouseService.UpdateWarehouse(&req))
	return
}
