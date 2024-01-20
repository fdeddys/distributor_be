package services

import (
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SupplierService struct {
}

// SaveDataSupplier Data Supplier
func (s SupplierService) SaveDataSupplier(suppliers *dbmodels.Supplier) models.ResponseSupplier {
	// var supplier dbmodels.Supplier

	if suppliers.ID < 1 {
		suppliers.Code = getSupplierCode()
	}
	suppliers.LastUpdate = time.Now()
	suppliers.LastUpdateBy = dto.CurrUser

	res := database.SaveSupplier(suppliers)
	fmt.Println("save : ", suppliers)

	return res
}

func getSupplierCode() string {
	db := database.GetDbCon()
	db.Debug().LogMode(true)

	prefix := "S"

	// err := db.Order(order).First(&models)
	var supplier []dbmodels.Supplier
	err := db.Model(&dbmodels.Supplier{}).Order("id desc").First(&supplier).Error
	// err := db.Model(&dbmodels.Brand{}).Where("id = 200000").Order("id desc").First(&brand).Error

	if err != nil {
		return prefix + "001"
	}
	if len(supplier) > 0 {
		// fmt.Printf("ini latest code nya : %s \n", brand[0].Code)
		woprefix := strings.TrimPrefix(supplier[0].Code, prefix)

		latestCode, err := strconv.Atoi(woprefix)
		if err != nil {
			fmt.Printf(err.Error())
			return prefix + "001"
		}
		// fmt.Printf("ini latest code nya : %d \n", latestCode)
		wpadding := fmt.Sprintf("%v%03s", prefix, strconv.Itoa(latestCode+1))
		// fmt.Printf("ini pake padding : %s \n", "B"+wpadding)
		return wpadding
	}
	return prefix + "001"

}

// Update Data Supplier
func (s SupplierService) UpdateDataSupplier(supplier *dbmodels.Supplier) models.Response {
	var data dbmodels.Supplier

	data.ID = supplier.ID
	data.Name = supplier.Name
	data.Code = supplier.Code
	data.Alamat = supplier.Alamat
	data.Kota = supplier.Kota
	data.Status = supplier.Status
	data.LastUpdate = time.Now()
	data.LastUpdateBy = dto.CurrUser
	data.PicName = supplier.PicName
	data.PicPhone = supplier.PicPhone
	data.BankAccountName = supplier.BankAccountName

	res := database.UpdateSupplier(&data)
	fmt.Println("update : ", res)

	return res
}

// Get Data Supplier Paging
func (s SupplierService) GetDataSupplierPaging(param dto.FilterPaging, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetSupplierPaging(param, offset, limit)

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

// GetBrandFilter ...
func (h SupplierService) GetSupplierByID(id int) models.ContentResponse {

	var res models.ContentResponse

	data, errCode, errDesc, err := database.GetSupplierByID(id)

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
