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

// SalesmanController ...
type SalesmanController struct {
	DB *gorm.DB
}

// SalesmanService ...
var SalesmanService = new(services.SalesmanService)

// SalesmanController ...
func (h *SalesmanController) GetSalesman(c *gin.Context) {

	res := SalesmanService.GetAllSalesman()

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetSalesmanFilter ...
func (h *SalesmanController) GetSalesmanFilter(c *gin.Context) {
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

	res = SalesmanService.GetSalesmanFilterPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetFilterBrand ...
func (h *SalesmanController) GetFilterSalesman(c *gin.Context) {
	res := models.ContentResponse{}

	id, errID := strconv.Atoi(c.Param("id"))
	if errID != nil {
		logs.Info("error", errID)
		// res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = SalesmanService.GetSalesmanFilter(id)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

//GetSalesmanLike ...
func (h *SalesmanController) GetSalesmanLike(c *gin.Context) {
	res := models.ContentResponse{}

	terms := c.Query("terms")

	if terms == "" {
		logs.Info("error", "can't found the salesman string")
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	// fmt.Sprintf("ini lookupstr = " + lookupstr)

	res = SalesmanService.GetSalesmanLike(terms)
	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// Save Salaeman ...
func (h *SalesmanController) SaveSalesman(c *gin.Context) {

	req := dbmodels.Salesman{}
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

	c.JSON(http.StatusOK, SalesmanService.SaveSalesman(&req))
	return
}

// UpdateBrand ...
func (h *SalesmanController) UpdateSalesman(c *gin.Context) {
	req := dbmodels.Salesman{}
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

	c.JSON(http.StatusOK, SalesmanService.UpdateSalesman(&req))
	return
}
