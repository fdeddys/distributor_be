package database

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetProductGroupDetails ...
func GetProductGroupDetails(id int) ([]dbmodels.ProductGroup, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var productGroup []dbmodels.ProductGroup
	err := db.Model(&dbmodels.ProductGroup{}).Where("id = ?", &id).First(&productGroup).Error

	if err != nil {
		return nil, constants.ERR_CODE_30, constants.ERR_CODE_30_MSG, err
	}
	// else {
	return productGroup, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
	// }
}

// GetProductGroupPaging ...
func GetProductGroupPaging(param dto.FilterProductGroup, offset int, limit int) ([]dbmodels.ProductGroup, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var productGroup []dbmodels.ProductGroup
	var total int
	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&productGroup).Error
		if err != nil {
			return productGroup, 0, err
		}
		return productGroup, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuery(db, offset, limit, &productGroup, param, "name", errQuery)
	go AsyncQueryCount(db, &total, param, &dbmodels.ProductGroup{}, "name", errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return productGroup, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("err-->", resErrCount)
		return productGroup, 0, resErrCount
	}

	return productGroup, total, nil
}

// UpdateProductGroup ...
func UpdateProductGroup(updatedProductGroup dbmodels.ProductGroup) models.NoContentResponse {
	var res models.NoContentResponse
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	db := GetDbCon()
	db.Debug().LogMode(true)

	var productGroup dbmodels.ProductGroup
	err := db.Model(&dbmodels.Brand{}).Where("id=?", &updatedProductGroup.ID).First(&productGroup).Error
	if err != nil {
		res.ErrCode = constants.ERR_CODE_30
		res.ErrDesc = constants.ERR_CODE_30_MSG + " " + err.Error()
		return res
	}

	productGroup.Name = updatedProductGroup.Name
	productGroup.Status = updatedProductGroup.Status
	productGroup.LastUpdateBy = dto.CurrUser
	productGroup.LastUpdate = time.Now()
	productGroup.Code = updatedProductGroup.Code

	err2 := db.Save(&productGroup)
	if err2 != nil {
		res.ErrCode = constants.ERR_CODE_30
		res.ErrDesc = constants.ERR_CODE_30_MSG + " " + err2.Error.Error()
	}

	return res
}

//SaveProductGroup ...
func SaveProductGroup(productGroup dbmodels.ProductGroup) models.NoContentResponse {
	var res models.NoContentResponse
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	db := GetDbCon()
	db.Debug().LogMode(true)
	if productGroup.ID < 1 {
		productGroup.Code = GenerateProductGroupCode()
	}

	productGroup.LastUpdateBy = dto.CurrUser
	productGroup.LastUpdate = time.Now()
	r := db.Save(&productGroup)
	if r.Error != nil {
		res.ErrCode = constants.ERR_CODE_30
		res.ErrDesc = constants.ERR_CODE_30_MSG + " " + r.Error.Error()
	}
	return res
}

func GenerateProductGroupCode() string {
	db := GetDbCon()
	db.Debug().LogMode(true)

	// err := db.Order(order).First(&models)
	var productGroup []dbmodels.ProductGroup
	err := db.Model(&dbmodels.ProductGroup{}).Order("id desc").First(&productGroup).Error
	// err := db.Model(&dbmodels.Brand{}).Where("id = 200000").Order("id desc").First(&brand).Error

	if err != nil {
		return "G001"
	}
	if len(productGroup) > 0 {
		// fmt.Printf("ini latest code nya : %s \n", brand[0].Code)
		woprefix := strings.TrimPrefix(productGroup[0].Code, "G")
		latestCode, err := strconv.Atoi(woprefix)
		if err != nil {
			fmt.Printf("error")
			return "G001"
		}
		// fmt.Printf("ini latest code nya : %d \n", latestCode)
		wpadding := fmt.Sprintf("%03s", strconv.Itoa(latestCode+1))
		// fmt.Printf("ini pake padding : %s \n", "B"+wpadding)
		return "G" + wpadding
	}
	return "G001"

}
