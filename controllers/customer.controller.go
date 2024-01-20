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
	"distribution-system-be/models/dbModels"
	"log"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CustomerController ...
type CustomerController struct {
	DB *gorm.DB
}

// CustomerService ...
var CustomerService = new(services.CustomerService)

// Save Data Customer
func (m *CustomerController) SaveDataCustomer(c *gin.Context) {
	CustomerReq := dbmodels.Customer{}
	res := models.Response{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &CustomerReq); err != nil {
		fmt.Println("Error, body Request")
		res.ErrCode = constants.ERR_CODE_03
		res.ErrDesc = constants.ERR_CODE_03_MSG
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, CustomerService.SaveDataCustomer(&CustomerReq))

	return
}

// List and Paging Customer
func (m *CustomerController) FilterDataCustomer(c *gin.Context) {
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
	res = CustomerService.GetDataCustomerPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	return
}

// Edit Data Customer
func (m *CustomerController) EditDataCustomer(c *gin.Context) {
	req := dbmodels.Customer{}
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
	c.JSON(http.StatusOK, CustomerService.UpdateDataCustomer(&req))
}

func (m *CustomerController) ListDataCustomerByName(c *gin.Context) {
	res := models.ContentResponse{}

	name := c.Query("search")
	if name == "" {
		logs.Info("error", "can't found the name string")
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	res = CustomerService.GetDataCustomerListByName(name)
	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// func (m *CustomerController) CheckOrderCustomerSupplier(c *gin.Context) {
// 	var checkSupplierReq dbmodels.CustomerCheckSupplier
// 	var res models.ResponseSFA
// 	var data string

// 	body := c.Request.Body
// 	dataBodyReq, _ := ioutil.ReadAll(body)

// 	if err := json.Unmarshal(dataBodyReq, &checkSupplierReq); err != nil {
// 		fmt.Println("Error, body Request")
// 		res.Data = ""
// 		res.Meta.Status = false
// 		res.Meta.Code = 400
// 		res.Meta.Message = "Terjadi Kesalahan"
// 		c.JSON(http.StatusBadRequest, res)
// 		c.Abort()
// 		return
// 	}

// 	checkSupplier := CustomerService.GetDataCheckOrder(checkSupplierReq.SupplierID, checkSupplierReq.CustomerID)

// 	if len(checkSupplier) > 0 {
// 		data = "Data tersedia"
// 	}else {
// 		data = "Data tidak tersedia"
// 	}

// 	res.Data = data
// 	res.Meta.Status = true
// 	res.Meta.Code = 200
// 	res.Meta.Message = "OK"

// 	c.JSON(http.StatusOK, res)
// 	return
// }
