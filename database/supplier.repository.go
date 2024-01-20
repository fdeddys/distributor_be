package database

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"log"
	_ "strconv"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// Get Data Supplier
func GetSupplierPaging(param dto.FilterPaging, offset int, limit int) ([]dbmodels.Supplier, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var supplier []dbmodels.Supplier
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&supplier).Error
		if err != nil {
			return supplier, 0, err
		}
		return supplier, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysSupplier(db, offset, limit, &supplier, param, errQuery)
	go AsyncQueryCountsSupplier(db, &total, param, errCount)
	fmt.Println(errQuery)
	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return supplier, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return supplier, 0, resErrCount
	}

	return supplier, total, nil
}

// AsyncQueryCountsSupplier ...
func AsyncQueryCountsSupplier(db *gorm.DB, total *int, param dto.FilterPaging, resChan chan error) {
	var searchName = "%"
	if strings.TrimSpace(param.Name) != "" {
		searchName += param.Name + "%"
	}

	var searchCode = "%"
	if strings.TrimSpace(param.Code) != "" {
		searchCode += param.Code + "%"
	}

	err := db.Model(&dbmodels.Supplier{}).Where("name ilike ? AND code ilike ?", searchName, searchCode).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysSupplier ...
func AsyncQuerysSupplier(db *gorm.DB, offset int, limit int, supplier *[]dbmodels.Supplier, param dto.FilterPaging, resChan chan error) {
	var searchName = "%"
	if strings.TrimSpace(param.Name) != "" {
		searchName += param.Name + "%"
	}

	var searchCode = "%"
	if strings.TrimSpace(param.Code) != "" {
		searchCode += param.Code + "%"
	}

	err := db.Preload("Bank").Order("id asc").Offset(offset).Limit(limit).Find(&supplier, "name ilike ? AND code ilike ?", searchName, searchCode).Error
	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// SaveSupplier Save
func SaveSupplier(supplier *dbmodels.Supplier) models.ResponseSupplier {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var res models.ResponseSupplier

	if r := db.Save(&supplier); r.Error != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Failed save data to DB"
		res.Code = ""
	}

	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG
	res.Code = supplier.Code

	return res
}

// Repository Update
func UpdateSupplier(supplier *dbmodels.Supplier) models.Response {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var res models.Response

	if r := db.Save(&supplier); r.Error != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Failed save data to DB"
	}

	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	return res
}

// Get Last Supplier
func GetLastSupplier() (dbmodels.Supplier, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var supplier dbmodels.Supplier
	var err error
	err = db.Order("code desc limit 1").Find(&supplier).Error
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

// GetListSupplier ...
func GetListSupplier() []dbmodels.Supplier {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var supplier []dbmodels.Supplier

	err := db.Find(&supplier).Error
	if err != nil {
		return supplier
	}

	return supplier

}

// GetSupplierByID ...
func GetSupplierByID(id int) (dbmodels.Supplier, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var supplier dbmodels.Supplier
	err := db.Preload("Bank").Model(&dbmodels.Supplier{}).Where("id = ?", &id).First(&supplier).Error

	if err != nil {
		return supplier, "02", "Error query data to DB", err
	}
	// else {
	return supplier, "00", "success", nil
	// }
}
