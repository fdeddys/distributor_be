package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CustomerService struct {
}

// Get Data Customer Paging
func (m CustomerService) GetDataCustomerPaging(param dto.FilterName, page int, limit int) models.ResponsePagination {
	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := database.GetCustomerPaging(param, offset, limit)
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

// SaveDataCustomer Data Customer
func (m CustomerService) SaveDataCustomer(Customer *dbmodels.Customer) models.Response {

	// if len(listCustomer) > 0 {
	// 	codeCustomer = strings.TrimPrefix(listCustomer[0].Code, string(listCustomer[0].Code[0]))
	// 	code, err := strconv.ParseInt(codeCustomer, 10, 64)
	// 	if err != nil {
	// 		var res models.Response
	// 		res.ErrCode = "05"
	// 		res.ErrCode = "Failed parse code Customer to integer"
	// 	}
	// 	code = code + 1
	// 	codeCustomer = fmt.Sprintf("%06d", code)
	// }else{
	// 	code = 1;
	// 	codeCustomer = fmt.Sprintf("%06d", code)
	// }

	if Customer.ID < 1 {
		Customer.Code = generateCustomerCode()
	}

	Customer.LastUpdate = time.Now()
	Customer.LastUpdateBy = dto.CurrUser

	res := database.SaveCustomer(Customer)
	fmt.Println("save : ", Customer)

	return res
}

func generateCustomerCode() string {

	var listCustomer dbmodels.Customer
	var code int64
	var codeCustomer string

	listCustomer, error := database.GetLastCustomer()
	code = 0

	if error != nil {
		code = 1
	} else {
		if listCustomer != (dbmodels.Customer{}) {
			if listCustomer.Code == "" {
				code = code + 1
			} else {
				codeCustomer = strings.TrimPrefix(listCustomer.Code, string(listCustomer.Code[0]))
				code, error = strconv.ParseInt(codeCustomer, 10, 64)
				code = code + 1
			}
		} else {
			code = 1
		}
	}
	codeCustomer = "C" + fmt.Sprintf("%07d", code)

	return codeCustomer
}

// UpdateDataCustomer Data Customer
func (m CustomerService) UpdateDataCustomer(Customer *dbmodels.Customer) models.Response {
	var data dbmodels.Customer
	data.ID = Customer.ID
	data.Name = Customer.Name
	data.Code = Customer.Code
	data.Top = Customer.Top
	data.Status = Customer.Status
	data.LastUpdateBy = dto.CurrUser
	data.LastUpdate = time.Now()

	// var res models.Response
	res := database.UpdateCustomer(&data)
	fmt.Println("update : ", res)

	return res
}

func (m CustomerService) GetDataCustomerListByName(name string) models.ContentResponse {
	var res models.ContentResponse

	data, err := database.GetListCustomerBySearch(name)

	if err != nil {
		res.ErrCode = "05"
		res.ErrDesc = "Failed load data"
		return res
	}

	res.ErrCode = "00"
	res.ErrDesc = constants.ERR_CODE_00_MSG
	res.Contents = data

	return res
}

func (m CustomerService) GetCustomerById(merchant_id int64) dbmodels.Customer {
	var res dbmodels.Customer
	res = database.GetCustomerById(merchant_id)
	return res
}

// func (m CustomerService) GetDataCheckOrder(supplier string, Customer string) []dbmodels.Order {
// 	res := database.GetOrderBySupplierAndCustomer(supplier, Customer)
// 	return res
// }
