package services

import (
	"distribution-system-be/database"
	"distribution-system-be/models"
)

// StockService ...
type StockService struct {
}

// GetDataStockMutationById ...
func (s StockService) GetDataStockByProductPage(productID int64, page, limit int) models.ResponsePagination {

	res := models.ResponsePagination{}
	// var err error
	offset := (page - 1) * limit
	stocks, totalData, err := database.GetStockByProductPage(productID, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = stocks
	res.Count = totalData
	res.Page = page
	res.Count = limit

	return res
}
