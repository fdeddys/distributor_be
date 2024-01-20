package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	cons "distribution-system-be/constants"
	models "distribution-system-be/models"
	dto "distribution-system-be/models/dto"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AuthController ...
type AuthController struct {
	DB *gorm.DB
}

// Login ...
func (h *AuthController) Login(c *gin.Context) {
	req := dto.LoginRequestDto{}
	res := models.Response{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = cons.ERR_CODE_03
		res.ErrDesc = cons.ERR_CODE_03_MSG
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, UserService.AuthLogin(&req))
}

// ChangePass ...
func (h *AuthController) ChangePass(c *gin.Context) {
	req := dto.ChangePassRequestDto{}
	res := models.Response{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = cons.ERR_CODE_03
		res.ErrDesc = cons.ERR_CODE_03_MSG
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	// if req.NewPassword != req.OldPassword {
	// 	res.ErrCode = "60"
	// 	res.ErrDesc = "New password and Last password not equal !"
	// 	c.JSON(http.StatusBadRequest, res)
	// 	c.Abort()
	// 	return
	// }

	c.JSON(http.StatusOK, UserService.ChangePass(&req))
}

// GetCurrPass ...
func (h *AuthController) GetCurrPass(c *gin.Context) {

	var res models.Response
	res.ErrCode = cons.ERR_CODE_00
	res.ErrDesc = dto.CurrUser

	c.JSON(http.StatusOK, res)
}
