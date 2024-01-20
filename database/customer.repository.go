package database

import (
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"distribution-system-be/constants"

	"github.com/jinzhu/gorm"
)

// Get Data Customer
func GetCustomerPaging(param dto.FilterName, offset int, limit int) ([]dbmodels.Customer, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var Customer []dbmodels.Customer
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&Customer).Error
		if err != nil {
			return Customer, 0, err
		}
		return Customer, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysCustomer(db, offset, limit, &Customer, param, errQuery)
	go AsyncQueryCountsCustomer(db, &total, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return Customer, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return Customer, 0, resErrCount
	}
	return Customer, total, nil
}

func AsyncQueryCountsCustomer(db *gorm.DB, total *int, param dto.FilterName, resChan chan error) {
	var searchName = "%"
	if strings.TrimSpace(param.Name) != "" {
		searchName = "%" + param.Name + "%"
	}
	// err := db.Table("issuer").Select("issuer.*, Customer.*").Joins("left join issuer on Customer.issuer_code = issuer.id").Model(&dbmodels.Customer{}).Where("Customer.name ilike ?", searchName).Count(&*total).Error
	err := db.Model(&dbmodels.Customer{}).Where("name ilike ?", searchName).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil

}

// AsyncQuerys ...
func AsyncQuerysCustomer(db *gorm.DB, offset int, limit int, Customer *[]dbmodels.Customer, param dto.FilterName, resChan chan error) {

	var searchName = "%"
	if strings.TrimSpace(param.Name) != "" {
		searchName = "%" + param.Name + "%"
	}

	err := db.Order("name ASC").Offset(offset).Limit(limit).Find(&Customer, "name ilike ?", searchName).Error
	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// Repository Save
func SaveCustomer(Customer *dbmodels.Customer) models.Response {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var res models.Response
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	if Customer.ID < 1 {
		Customer.Code = GenerateCustomerCode()
	}
	r := db.Save(&Customer)
	if r.Error != nil {
		res.ErrCode = constants.ERR_CODE_30
		res.ErrDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}

	return res
}

// UpdateCustomer Update
func UpdateCustomer(Customer *dbmodels.Customer) models.Response {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var res models.Response
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	if r := db.Save(&Customer); r.Error != nil {
		res.ErrCode = constants.ERR_CODE_30
		res.ErrDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}

	return res
}

func GetListCustomer() ([]dbmodels.Customer, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var Customer []dbmodels.Customer
	var err error

	err = db.Find(&Customer).Error
	if err != nil {
		return Customer, err
	}
	return Customer, nil
}

// Get Last Customer
func GetLastCustomer() (dbmodels.Customer, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var Customer dbmodels.Customer
	var err error

	err = db.Order("code desc limit 1").Find(&Customer).Error
	if err != nil {
		return Customer, err
	}
	return Customer, nil
}

func GetListCustomerBySearch(name string) ([]dbmodels.Customer, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var Customer []dbmodels.Customer
	var err error

	err = db.Where("name ilike ? ", "%"+name+"%").Find(&Customer).Error
	if err != nil {
		return Customer, err
	}
	return Customer, nil
}

// Customer check supplier
func GetOrderBySupplierAndCustomer(supplier string, Customer string) []dbmodels.SalesOrder {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var order []dbmodels.SalesOrder
	var err error

	err = db.Where("supplier_code ilike ? and Customer_code ilike ?", supplier, Customer).Find(&order).Error

	if err != nil {
		return order
	}

	return order
}

// Customer by id
func GetCustomerById(Customer_id int64) dbmodels.Customer {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var Customer dbmodels.Customer
	var err error

	err = db.Where("id = ?", Customer_id).Find(&Customer).Error

	if err != nil {
		return Customer
	}
	return Customer
}

func FindCustomerByPhone(phoneNumb string) dbmodels.Customer {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var Customer dbmodels.Customer
	var err error

	err = db.Order("id desc limit 1").Where("phone_numb = ?", phoneNumb).Find(&Customer).Error

	if err == nil {
		return Customer
	}
	return dbmodels.Customer{}
}

// GenerateCustomerCode ...
func GenerateCustomerCode() string {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var customer []dbmodels.Customer
	err := db.Model(&dbmodels.Customer{}).Order("id desc").First(&customer).Error

	if err != nil {
		return "C0000001"
	}
	if len(customer) > 0 {
		// fmt.Printf("ini latest code nya : %s \n", brand[0].Code)
		prefix := strings.TrimPrefix(customer[0].Code, "C")
		latestCode, err := strconv.Atoi(prefix)
		if err != nil {
			fmt.Printf("error")
			return "C0000001"
		}
		wpadding := fmt.Sprintf("%06s", strconv.Itoa(latestCode+1))
		return "C" + wpadding
	}
	return "C0000001"

}
