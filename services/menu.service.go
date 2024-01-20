package services

import (
	repository "distribution-system-be/database"
	dbmodels "distribution-system-be/models/dbModels"
)

// MenuService ...
type MenuService struct {
}

// GetMenuByUser ...
func (h MenuService) GetMenuByUser(user string) []dbmodels.Menu {
	var res []dbmodels.Menu
	// var err error
	res, _ = repository.GetUserMenus(user)

	return res
}
