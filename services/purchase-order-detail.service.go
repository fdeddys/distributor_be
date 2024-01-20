package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
)

// PurchaseOrderDetailService ...
type PurchaseOrderDetailService struct {
}

// GetDataPurchaseOrderDetailPage ...
func (r PurchaseOrderDetailService) GetDataPurchaseOrderDetailPage(param dto.FilterPurchaseOrderDetail, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetPurchaseOrderDetailPage(param, offset, limit)

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

// SavePurchaseOrderDetail ...
func (r PurchaseOrderDetailService) SavePurchaseOrderDetail(purchaseOrderDetail *dbmodels.PurchaseOrderDetail) (errCode string, errDesc string) {

	if _, err := database.GetPurchaseOrderByPurchaseOrderID(purchaseOrderDetail.PurchaseOrderID); err != nil {
		return "99", err.Error()
	}

	curProd, _, _ := database.FindProductByID(purchaseOrderDetail.ProductID)
	purchaseOrderDetail.Qty = purchaseOrderDetail.PoQty * purchaseOrderDetail.PoUOMQty
	// purchaseOrderDetail.UomID = purchaseOrderDetail.PoUomID
	purchaseOrderDetail.UomID = curProd.SmallUomID
	purchaseOrderDetail.Price = (purchaseOrderDetail.PoPrice / float32(purchaseOrderDetail.PoUOMQty))
	if err, errDesc := database.SavePurchaseOrderDetail(purchaseOrderDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	CalculateTotalPO(purchaseOrderDetail.PurchaseOrderID)
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// DeletePurchaseOrderDetailByID ...
func (r PurchaseOrderDetailService) DeletePurchaseOrderDetailByID(purchaseOrderDetailID int64) (errCode string, errDesc string) {

	po, _ := database.GetPurchaseOrderByPurchaseOrderDetailID(purchaseOrderDetailID)
	if err, errDesc := database.DeletePurchaseOrderDetailById(purchaseOrderDetailID); err != constants.ERR_CODE_00 {
		return err, errDesc
	}

	CalculateTotalPO(po.ID)
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetDataPurchaseOrderDetailPage ...
func (r PurchaseOrderDetailService) GetLastPricePurchaseOrderDetail(productId int64) models.ResponseCheckPrice {
	var res models.ResponseCheckPrice

	data := database.GetLastPricePurchaseOrderDetail(productId)
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG
	res.Price = data.Price
	res.Disc1 = data.Disc1
	res.Disc2 = data.Disc2
	res.Hpp = data.Hpp
	return res
}

func (r PurchaseOrderDetailService) UpdateDetail(purchaseOrderDetail *dbmodels.PurchaseOrderDetail) (errCode string, errDesc string) {

	// for _, receiveDetail := range *receiveDetails {
	// fmt.Println(receiveDetail.ID, "  ", receiveDetail.Qty, " ", receiveDetail.Price, " ", receiveDetail.Disc1)
	if err, errDesc := database.UpdatePODetail(*purchaseOrderDetail); err != constants.ERR_CODE_00 {
		return err, errDesc
	}
	// }

	// if err, errDesc := database.SaveReceiveDetail(receiveDetail); err != constants.ERR_CODE_00 {
	// 	return err, errDesc
	// }

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
