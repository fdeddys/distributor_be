package dto

// CurrUser ...
var CurrUser string
var CurrUserId int64

// FilterDto ...
type FilterDto struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// FilterUser ...
type FilterUser struct {
	Username string `json:"username"`
}

// FilterMenu ...
type FilterMenu struct {
	MenuName string `json:"menuname"`
}

// FilterBrand ...
type FilterBrand struct {
	Name string `json:"name"`
}

// FilterProduct ...
type FilterProduct struct {
	Name        string `json:"name"`
	WarehouseID int64  `json:"warehouseId"`
	Composition string `json:"composition"`
}

// FilterProductGroup ...
type FilterProductGroup struct {
	Name string `json:"name"`
}

// LoginRequestDto ...
type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePassRequestDto struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// LoginResponseDto ...
type LoginResponseDto struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Token   string `json:"token"`
}

// FilterPaging ...
type FilterPaging struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// FilterMerchant ...
type FilterName struct {
	Name string `json:"name"`
}

// FilterLookup ...
type FilterLookup struct {
	Name string `json:"group_name"`
}

type FilterSupplierGroup struct {
	Name string `json:"name"`
}

type FilterFollowUpOrder struct {
	OrderNumber string `json:"order_no"`
}

// NoContentResponse ...
type NoContentResponse struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
}
