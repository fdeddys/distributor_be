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

// BrandController ...
type BrandController struct {
	DB *gorm.DB
}

// BrandService ...
var BrandService = new(services.BrandService)

// GetBrand ...
func (h *BrandController) GetBrand(c *gin.Context) {
	req := dto.FilterBrand{}
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

	res = BrandService.GetBrandFilterPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetFilterBrand ...
func (h *BrandController) GetFilterBrand(c *gin.Context) {
	res := models.ContentResponse{}

	id, errID := strconv.Atoi(c.Param("id"))
	if errID != nil {
		logs.Info("error", errID)
		// res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = BrandService.GetBrandFilter(id)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

//GetBrandLike ...
func (h *BrandController) GetBrandLike(c *gin.Context) {
	res := models.ContentResponse{}

	brandterms := c.Query("terms")

	if brandterms == "" {
		logs.Info("error", "can't found the brand string")
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	// fmt.Sprintf("ini lookupstr = " + lookupstr)

	res = BrandService.GetBrandLike(brandterms)
	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// SaveBrand ...
func (h *BrandController) SaveBrand(c *gin.Context) {

	req := dbmodels.Brand{}
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

	c.JSON(http.StatusOK, BrandService.SaveBrand(&req))
	return
}

// UpdateBrand ...
func (h *BrandController) UpdateBrand(c *gin.Context) {
	req := dbmodels.Brand{}
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

	c.JSON(http.StatusOK, BrandService.UpdateBrand(&req))
	return
}
