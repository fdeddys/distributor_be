package controllers

import (
	"distribution-system-be/services"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models"

)

// OrderController ...
type ParameterController struct {
	DB *gorm.DB
}

// ParameterService ...
var ParameterService = new(services.ParameterService)

// GetByName ...
func (s *ParameterController) GetByName(c *gin.Context) {

	// res := dbmodels.Parameter{}

	paramName := c.Param("param-name")
	res := ParameterService.GetByName(paramName)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// GetAll ...
func (s *ParameterController) GetAll(c *gin.Context) {

	// res := dbmodels.Parameter{}

	res := ParameterService.GetAll()

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// GetAll ...
func (s *ParameterController) UpdateParam(c *gin.Context) {

	req := dbmodels.Parameter{}
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

	c.JSON(http.StatusOK, ParameterService.UpdateParam(&req))
	return

}
