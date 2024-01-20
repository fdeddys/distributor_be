package controllers

import (
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/services"
	"distribution-system-be/utils/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"distribution-system-be/constants"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserController ...
type UserController struct {
	DB *gorm.DB
}

// UserService ...
var UserService = new(services.UserService)

// GetUser ...
func (h *UserController) GetUser(c *gin.Context) {
	req := dto.FilterUser{}
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

	// log.Println("page->", page, "count->", count)

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	// temp, _ := json.Marshal(req)

	// log.Println("req-->", string(temp))

	res = UserService.GetUserFilterPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return

}

// SaveDataUser ...
func (h *UserController) SaveDataUser(c *gin.Context) {
	fmt.Println("entering the save data user ")
	req := dbmodels.User{}
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

	c.JSON(http.StatusOK, UserService.SaveUser(&req))
}

// UpdateUser ...
func (h *UserController) UpdateUser(c *gin.Context) {
	req := dbmodels.User{}
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

	c.JSON(http.StatusOK, UserService.UpdateUser(&req))
	return
}

// ResetUser ...
func (h *UserController) ResetUser(c *gin.Context) {
	res := models.ContentResponse{}

	idUser := util.Atoi64(c.Param("iduser"))
	if idUser == 0 {
		logs.Info("error id")
		res.ErrCode = constants.ERR_CODE_03
		res.ErrCode = "Error get user id"
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, UserService.ResetUser(idUser))
	return
}

// func (p *UserController) EditDataUser(c *gin.Context) {
// 	req := dbmodels.UsersTransaction{}
// 	res := models.Response{}
// 	body := c.Request.Body
// 	dataBodyReq, _ := ioutil.ReadAll(body)

// 	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
// 		fmt.Println("Error, body Request ")
// 		res.ErrCode = "03"
// 		res.ErrDesc = "Error, unmarshall body Request"
// 		c.JSON(http.StatusBadRequest, res)
// 		c.Abort()
// 		return
// 	}
// 	fmt.Println("Edit data")
// 	c.JSON(http.StatusOK, UserService.UpdateDataUser(&req))
// }

// func (p *UserController) ChangePassword(c *gin.Context) {
// 	req := dbmodels.ChangePassword{}
// 	res := models.Response{}
// 	body := c.Request.Body
// 	dataBodyReq, _ := ioutil.ReadAll(body)

// 	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
// 		fmt.Println("Error, body Request ")
// 		res.ErrCode = "03"
// 		res.ErrDesc = "Error, unmarshall body Request"
// 		c.JSON(http.StatusBadRequest, res)
// 		c.Abort()
// 		return
// 	}
// 	fmt.Println("Edit data Password")
// 	c.JSON(http.StatusOK, UserService.UpdateChangePassword(&req))
// }
