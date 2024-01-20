package controllers

import (
	"distribution-system-be/models"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

// SupplierController ...
type SupplierController struct {
	DB *gorm.DB
}

// SupplierService ...
var SupplierService = new(services.SupplierService)

var ProductServices = new(services.ProductService)
var CustomerServices = new(services.CustomerService)

/* ------------------------------------------- Begin Supplier ----------------------------------------------- */

func (s *SupplierController) SaveDataSupplier(c *gin.Context) {
	SupplierReq := dbmodels.Supplier{}
	res := models.ResponseSupplier{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &SupplierReq); err != nil {
		fmt.Println("Error, body Request")
		res.ErrCode = constants.ERR_CODE_03
		res.ErrDesc = constants.ERR_CODE_03_MSG
		res.Code = ""
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, SupplierService.SaveDataSupplier(&SupplierReq))

	return
}

func (s *SupplierController) FilterDataSupplier(c *gin.Context) {
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
	res = SupplierService.GetDataSupplierPaging(req, page, count)

	c.JSON(http.StatusOK, res)

	return
}

func (s *SupplierController) EditDataSupplier(c *gin.Context) {
	req := dbmodels.Supplier{}
	res := models.Response{}

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

	fmt.Println("Edit data")
	c.JSON(http.StatusOK, SupplierService.UpdateDataSupplier(&req))
}

/* ----------------------------------- End Supplier --------------------------------------------------------- */

/* --------------------------------- Begin Supplier Merchant ------------------------------------------------ */

func (s *SupplierController) FilterDataSupplierMerchant(c *gin.Context) {
	req := dto.FilterSupplierMerchantDto{}
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

	supplier_id, errSupplierId := strconv.Atoi(c.Param("supplier_id"))
	if errSupplierId != nil {
		logs.Info("error", errPage)
		res.Error = errSupplierId.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	log.Println("page->", page, "count->", count, "supplier_id->", supplier_id)

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
	log.Println("searchCode-->", string(temp))

	c.JSON(http.StatusOK, res)

	return
}

// GetFilterBrand ...
func (s *SupplierController) FilterByID(c *gin.Context) {
	res := models.ContentResponse{}

	id, errID := strconv.Atoi(c.Param("id"))
	if errID != nil {
		logs.Info("error", errID)
		// res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = SupplierService.GetSupplierByID(id)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}
