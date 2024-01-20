package controllers

import (
	"distribution-system-be/models"
	"distribution-system-be/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//LookupGroupController ...
type LookupGroupController struct {
	DB *gorm.DB
}

//LookupService ...
var lookupGroupService = new(services.LookupGroupService)

//GetLookupGroup ...
func (h *LookupGroupController) GetLookupGroup(c *gin.Context) {
	res := models.ContentResponse{}

	res = lookupGroupService.GetLookupGroup()
	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}
