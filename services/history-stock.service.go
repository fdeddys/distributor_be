package services

import (
	"distribution-system-be/database"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"time"
)

// HistoryStockService ...
type HistoryStockService struct {
}

// SaveHistoryStock ...
func (h HistoryStockService) SaveHistoryStock(history *dbmodels.HistoryStock) models.NoContentResponse {
	history.LastUpdate = time.Now()
	history.LastUpdateBy = dto.CurrUser

	res := repository.SaveHistory(*history)
	return res
}

// GetHistoryStockPage
func (h HistoryStockService) GetHistoryStockPage(param dto.FilterHistoryStock, page, limit int) models.ResponsePagination {

	res := models.ResponsePagination{}
	offset := (page - 1) * limit
	historyStocks, totalData, err := database.GetHistoryPage(param, offset, limit)
	for idx, historyStock := range historyStocks {
		// find uom
		prod, _, _ := database.FindProductByCode(historyStock.Code)
		historyStocks[idx].Satuan = prod.SmallUom.Name
	}

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = historyStocks
	res.TotalRow = totalData
	res.Page = page
	res.Count = limit

	return res
}
