package database

import (
	constants "distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"
	"log"
	"strconv"
	"sync"
)

// UpdateBrand ...
func UpdateBrand(updatedbrand dbmodels.Brand) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	var brand dbmodels.Brand
	err := db.Model(&dbmodels.Brand{}).Where("id=?", &updatedbrand.ID).First(&brand).Error
	if err != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error select data to DB"
		return res
	}

	brand.Name = updatedbrand.Name
	brand.Status = updatedbrand.Status
	brand.LastUpdateBy = updatedbrand.LastUpdateBy
	brand.LastUpdate = updatedbrand.LastUpdate
	brand.Code = updatedbrand.Code

	err2 := db.Save(&brand)
	if err2 != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error update data to DB"
		return res
	}

	res.ErrCode = "00"
	res.ErrDesc = "Success"

	return res
}

//GetBrandLike ...
func GetBrandLike(brandTerms string) ([]dbmodels.Brand, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var brand []dbmodels.Brand
	err := db.Model(&dbmodels.Brand{}).Where("name ~* ?", &brandTerms).Find(&brand).Error

	if err != nil {
		return nil, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG, err
	}
	return brand, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
}

//SaveBrand ...
func SaveBrand(brand dbmodels.Brand) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	if brand.ID < 1 {
		brand.Code = GenerateBrandCode()
	}

	if r := db.Save(&brand); r.Error != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error save data to DB"
		return res
	}

	res.ErrCode = "00"
	res.ErrDesc = "Success"
	return res
}

// GetBrandFilter ...
func GetBrandFilter(id int) ([]dbmodels.Brand, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var brand []dbmodels.Brand
	err := db.Model(&dbmodels.Brand{}).Where("id = ?", &id).First(&brand).Error

	if err != nil {
		return nil, "02", "Error query data to DB", err
	}
	// else {
	return brand, "00", "success", nil
	// }
}

// GetBrand ...
func GetBrand(param dto.FilterBrand, offset int, limit int) ([]dbmodels.Brand, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var brand []dbmodels.Brand
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&brand).Error
		if err != nil {
			return brand, 0, err
		}
		return brand, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQueryCount(db, &total, param, &dbmodels.Brand{}, "name", errCount)
	// if limit == 0 {
	// 	limit = total
	// }
	go AsyncQuery(db, offset, limit, &brand, param, "name", errQuery)

	resErrCount := <-errCount
	resErrQuery := <-errQuery

	wg.Done()
	// wg.Done()

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return brand, 0, resErrCount
	}

	if resErrQuery != nil {
		return brand, 0, resErrQuery
	}

	return brand, total, nil
}

// func GenerateBrandCode() string {
// 	db := GetDbCon()
// 	db.Debug().LogMode(true)

// 	// err := db.Order(order).First(&models)
// 	var brand []dbmodels.Brand
// 	err := db.Model(&dbmodels.Brand{}).Order("id desc").First(&brand).Error
// 	// err := db.Model(&dbmodels.Brand{}).Where("id = 200000").Order("id desc").First(&brand).Error

// 	if err != nil {
// 		return "B001"
// 	}
// 	if len(brand) > 0 {
// 		// fmt.Printf("ini latest code nya : %s \n", brand[0].Code)
// 		woprefix := strings.TrimPrefix(brand[0].Code, "")
// 		latestCode, err := strconv.Atoi(woprefix)
// 		if err != nil {
// 			fmt.Printf("error")
// 			return "B001"
// 		}
// 		// fmt.Printf("ini latest code nya : %d \n", latestCode)
// 		wpadding := fmt.Sprintf("%03s", strconv.Itoa(latestCode+1))
// 		// fmt.Printf("ini pake padding : %s \n", "B"+wpadding)
// 		return wpadding
// 	}
// 	return "B001"

// }

func GenerateBrandCode() string {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var brands []dbmodels.Brand
	err := db.Model(&dbmodels.Brand{}).Order("id desc").Find(&brands).Error
	header := "B"
	defaultCode := "B001"

	if err != nil {
		return defaultCode
	}
	if len(brands) > 0 {
		// fmt.Printf("ini latest code nya : %s \n", salesman[0].Code)
		code := brands[0].Code
		runes := []rune(code)
		latestNumb := string(runes[1:len(code)])
		fmt.Println("latest numb-", latestNumb)
		// woprefix := strings.TrimPrefix(salesman[0].Code, "")
		latestCode, err := strconv.Atoi(latestNumb)
		if err != nil {
			fmt.Println("error =>", err.Error())
			return defaultCode
		}
		// fmt.Printf("ini latest code nya : %d \n", latestCode)
		wpadding := fmt.Sprintf("%v%03s", header, strconv.Itoa(latestCode+1))
		// fmt.Printf("ini pake padding : %s \n", "B"+wpadding)
		return wpadding
	}
	return defaultCode

}

// // AsyncQueryCount ...
// func AsyncQueryCount(db *gorm.DB, total *int, param interface{}, models interface{}, fieldLookup string, resChan chan error) {
// 	varInterface := reflect.ValueOf(param)
// 	strQuery := varInterface.Field(0).Interface().(string)

// 	var criteriaName = ""
// 	if strings.TrimSpace(strQuery) != "" {
// 		criteriaName = strQuery
// 	}

// 	err := db.Model(models).Where(fieldLookup+" ~* ?", criteriaName).Count(&*total).Error

// 	if err != nil {
// 		resChan <- err
// 	}
// 	resChan <- nil
// }

// // AsyncQuery ...
// func AsyncQuery(db *gorm.DB, offset int, limit int, modelWVal interface{}, param interface{}, fieldLookup string, isProduct bool, resChan chan error) {
// 	modelsDump := reflect.ValueOf(modelWVal).Interface()
// 	paramDump := reflect.ValueOf(param)
// 	strQuery := paramDump.Field(0).Interface().(string)
// 	var criteriaName = ""
// 	if strings.TrimSpace(strQuery) != "" {
// 		criteriaName = strQuery //+ criteriaBrandName
// 	}

// 	var err error
// 	if isProduct {
// 		err = db.Preload("ProductGroup").Preload("Brand").Order("name ASC").Offset(offset).Limit(limit).Find(modelsDump, fieldLookup+" ~* ?", criteriaName).Error
// 	} else {
// 		err = db.Order("name ASC").Offset(offset).Limit(limit).Find(modelsDump, fieldLookup+" ~* ?", criteriaName).Error
// 	}

// 	if err != nil {
// 		resChan <- err
// 	}
// 	resChan <- nil
// }
