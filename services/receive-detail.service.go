package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
)

// ReceiveDetailService ...
type ReceiveDetailService struct {
}

// GetDataReceiveDetailPage ...
func (r ReceiveDetailService) GetDataReceiveDetailPage(param dto.FilterReceiveDetail, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReceiveDetailPage(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = limit

	return res
}

// SaveReceiveDetail ...
func (r ReceiveDetailService) SaveReceiveDetail(receiveDetail *dbmodels.ReceiveDetail) (errCode string, errDesc string) {

	if _, err := database.GetReceiveByReceiveID(receiveDetail.ReceiveID); err != nil {
		return "99", err.Error()
	}

	if err, errDesc := database.SaveReceiveDetail(receiveDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// SaveReceiveDetail ...
func (r ReceiveDetailService) UpdateReceiveDetail(receiveDetails *[]dbmodels.ReceiveDetail) (errCode string, errDesc string) {

	for _, receiveDetail := range *receiveDetails {
		// fmt.Println(receiveDetail.ID, "  ", receiveDetail.Qty, " ", receiveDetail.Price, " ", receiveDetail.Disc1)
		if err, errDesc := database.UpdateReceiveDetail(receiveDetail.ID, receiveDetail.Qty, receiveDetail.Price, receiveDetail.Disc1, receiveDetail.Disc2, receiveDetail.BatchNo, receiveDetail.Ed); err != constants.ERR_CODE_00 {
			return err, errDesc
		}
	}

	// if err, errDesc := database.SaveReceiveDetail(receiveDetail); err != constants.ERR_CODE_00 {
	// 	return err, errDesc
	// }

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeleteReceiveDetailByID ...
func (r ReceiveDetailService) DeleteReceiveDetailByID(receiveDetailID int64) (errCode string, errDesc string) {

	if err, errDesc := database.DeleteReceiveDetailById(receiveDetailID); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func (r ReceiveDetailService) DeleteReceiveDetailByIDMultiple(receiveDetailIDs []int64) (errCode string, errDesc string) {

	for _, id := range receiveDetailIDs {
		fmt.Println("delete by id " , id)
		database.DeleteReceiveDetailById(id)
	}


	// if err, errDesc := database.DeleteReceiveDetailById(receiveDetailID); err != constants.ERR_CODE_00 {
	// 	return err, errDesc
	// }

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetDataPurchaseOrderDetailPage ...
func (r ReceiveDetailService) GetDataBatchExpired(param dto.FilterBatchExpired, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetReceiveOrderDetailBatchExpiredPage(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	// add stock
	for idx, receiveDetail := range data {
		
		// res := models.ResponsePagination{}
		// offset := (page - 1) * limit
		stocks, _, _ := database.GetStockByProductPage(receiveDetail.ProductID, 0, 5)
		if len(stocks)>0 {
			data[idx].QtyWh = stocks[0].Qty
			// fmt.Println("QtyWh = ", receiveDetail.QtyWh )
		}

		
	}
	// fmt.Println("data = ", data )

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = limit

	return res
}

// GetDataPriceProductlPage ...
func (r ReceiveDetailService) GetDataPriceProduct(productId int64) models.ResponseReceiveCheckPrice {

	var res models.ResponseReceiveCheckPrice
	data := database.GetDataPriceProduct(productId)

	fmt.Println("service : ", data)
	res.Price = data.Price

	return res
}
