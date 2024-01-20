package services

import (
	"distribution-system-be/constants"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"time"
	"fmt"
)

// ProductService ...
type ProductService struct {
}

// GetProductFilterPaging ...
func (h ProductService) GetProductFilterPaging(param dto.FilterProduct, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	// if validateMenu()==false{
	// 	res.Contents = nil
	// 	return res
	// }
	offset := (page - 1) * limit
	data, totalData, err := repository.GetProductListPaging(param, offset, limit, false)

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

// GetProductFilterPagingAllStatus ...
func (h ProductService) GetProductFilterPagingAllStatus(param dto.FilterProduct, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	// if validateMenu()==false{
	// 	res.Contents = nil
	// 	return res
	// }
	offset := (page - 1) * limit
	data, totalData, err := repository.GetProductListPaging(param, offset, limit, true)

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

func (h ProductService) SearchProduct(param dto.FilterProduct, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, err := repository.SearchProduct(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = 0
	res.Page = page
	res.Count = len(data)

	return res
}

// GetProductDetails ...
func (h ProductService) GetProductDetails(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetProductDetails(id)

	if err != nil {
		res.Contents = nil
		res.ErrCode = "02"
		res.ErrDesc = "Error query data to DB"
		return res
	}

	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	return res
}

// SaveProduct ...
func (h ProductService) SaveProduct(product *dbmodels.Product) models.ContentResponse {
	if validateMenu()==false{
		var resInvalid models.ContentResponse
		resInvalid.ErrCode="99"
		resInvalid.ErrDesc="Anauthorized, please relogin !!"
		resInvalid.Contents = nil
		return resInvalid
	}
	
	product.LastUpdate = time.Now()
	product.LastUpdateBy = dto.CurrUser

	// var res models.ResponseSave
	res := repository.SaveProduct(*product)

	return res
}

// UpdateProduct ...
func (h ProductService) UpdateProduct(product *dbmodels.Product) models.NoContentResponse {
	if validateMenu()==false{
		var resInvalid models.NoContentResponse
		resInvalid.ErrCode="99"
		resInvalid.ErrDesc="Anauthorized, please relogin !!"
		return resInvalid
	}
	
	product.LastUpdate = time.Now()
	product.LastUpdateBy = dto.CurrUser

	res := repository.UpdateProduct(*product)

	return res
}

func (h ProductService) ProductList() []dbmodels.Product {
	res := repository.ProductList()
	return res
}

// GetProductLike ...
func (h ProductService) GetProductLike(productterms string) models.ContentResponse {
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetProductLike(productterms)

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res
	}

	res.Contents = data
	res.ErrCode = errCode
	res.ErrDesc = errDesc

	return res
}

func validateMenu() bool {

	fmt.Println("Username : " + dto.CurrUser)
	result:= repository.ValidateMenuByUserName(dto.CurrUser, "product")

	return result

}
