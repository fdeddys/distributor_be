package services

import (
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"time"
)

// RoleService ...
type RoleService struct {
}

// GetRoleFilterPaging ...
func (h RoleService) GetRoleFilterPaging(param dto.FilterName, page int, limit int) models.ResponsePagination {

	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetRole(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = len(data)

	return res
}

// SaveRole ...
func (h RoleService) SaveRole(role *dbmodels.Role) models.NoContentResponse {

	role.LastUpdate = time.Now()
	role.LastUpdateBy = dto.CurrUser
	res, newID := repository.SaveRole(*role)
	repository.CreateRoleMenuByRoleId(newID)
	return res
}

// UpdateRole ...
func (h RoleService) UpdateRole(role *dbmodels.Role) models.NoContentResponse {

	res := repository.UpdateRole(*role)

	return res
}
