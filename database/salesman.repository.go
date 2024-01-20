package database

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"log"
	"strconv"
	"sync"
)

// GetSales ...
func GetAllSalesman() ([]dbmodels.Salesman, string, string) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var salesman []dbmodels.Salesman
	err := db.Where("status = ?", 1).Find(&salesman).Error
	if err != nil {
		return nil, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG + "  " + err.Error()
	}
	return salesman, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// UpdateSalesman ...
func UpdateSalesman(updatedSalesman dbmodels.Salesman) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	var salesman dbmodels.Salesman
	err := db.Model(&dbmodels.Salesman{}).Where("id=?", &updatedSalesman.ID).First(&salesman).Error
	if err != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error select data to DB"
	}

	salesman.Name = updatedSalesman.Name
	salesman.Status = updatedSalesman.Status
	salesman.LastUpdateBy = updatedSalesman.LastUpdateBy
	salesman.LastUpdate = updatedSalesman.LastUpdate
	salesman.Code = updatedSalesman.Code
	salesman.Description = updatedSalesman.Description

	err2 := db.Save(&salesman)
	if err2 != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error update data to DB"
	}

	res.ErrCode = "00"
	res.ErrDesc = "Success"

	return res
}

//GetBrandLike ...
func GetSalesmanLike(salesmanTerms string) ([]dbmodels.Salesman, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var salesman []dbmodels.Salesman
	err := db.Model(&dbmodels.Salesman{}).Where("name ~* ?", &salesmanTerms).Find(&salesman).Error

	if err != nil {
		return nil, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG, err
	}
	return salesman, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
}

//SaveSalesman ...
func SaveSalesman(salesman dbmodels.Salesman) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	if salesman.ID < 1 {
		salesman.Code = GenerateSalesmanCode()
	}

	if r := db.Save(&salesman); r.Error != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error save data to DB"
	}

	res.ErrCode = "00"
	res.ErrDesc = "Success"
	return res
}

// GetSalesmanFilter ...
func GetSalesmanFilter(id int) ([]dbmodels.Salesman, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var salesman []dbmodels.Salesman
	err := db.Model(&dbmodels.Salesman{}).Where("id = ?", &id).First(&salesman).Error

	if err != nil {
		return nil, "02", "Error query data to DB", err
	}
	return salesman, "00", "success", nil
}

// GetSalesman ...
func GetSalesman(param dto.FilterPaging, offset int, limit int) ([]dbmodels.Salesman, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var salesman []dbmodels.Salesman
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&salesman).Error
		if err != nil {
			return salesman, 0, err
		}
		return salesman, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQueryCount(db, &total, param, &dbmodels.Salesman{}, "name", errCount)
	go AsyncQuery(db, offset, limit, &salesman, param, "name", errQuery)

	resErrCount := <-errCount
	resErrQuery := <-errQuery

	wg.Done()
	// wg.Done()

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return salesman, 0, resErrCount
	}

	if resErrQuery != nil {
		return salesman, 0, resErrQuery
	}

	return salesman, total, nil
}

func GenerateSalesmanCode() string {
	db := GetDbCon()
	db.Debug().LogMode(true)
	header := "S"
	defaultCode := "S001"

	var salesmans []dbmodels.Salesman
	err := db.Model(&dbmodels.Salesman{}).Order("id desc").Find(&salesmans).Error

	if err != nil {
		return defaultCode
	}
	if len(salesmans) > 0 {
		// fmt.Printf("ini latest code nya : %s \n", salesman[0].Code)
		code := salesmans[0].Code
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
