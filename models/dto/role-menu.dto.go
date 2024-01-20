package dto

// RoleMenuDto ...
type RoleMenuDto struct {
	MenuID          int    `json:"menuId" gorm:"column:menuid"`
	MenuDescription string `json:"menuDescription" gorm:"column:menudescription"`
	Status          int    `json:"status" gorm:"column:status"`
	ParentId        int    `json:"parentId" gorm:"column:parentid"`
}

// ContentResponse ...
type ContentResponse struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Data    string `json:"data"`
}

// RequestUpdateRole ...
type RequestUpdateRole struct {
	RoleID int `json:"roleId"`
	MenuID int `json:"menuId"`
	Status int `json:"status"`
}
