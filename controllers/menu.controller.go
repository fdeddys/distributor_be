package controllers

import (
	"net/http"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/services"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// MenuController ...
type MenuController struct {
	DB *gorm.DB
}

// MenuService ...
var MenuService = new(services.MenuService)

// GetMenuByUser ...
func (h *MenuController) GetMenuByUser(c *gin.Context) {
	res := []dbmodels.Menu{}

	res = MenuService.GetMenuByUser(dto.CurrUser)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}
