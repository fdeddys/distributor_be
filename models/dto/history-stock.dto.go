package dto

type FilterHistoryStock struct {
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	ProductCode string `json:"productCode"`
	WarehouseID int64  `json:"warehouseId"`
}
