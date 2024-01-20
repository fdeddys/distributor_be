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
)

//GetLookupByGroup ...
func GetLookupByGroup(lookupstr string) ([]dbmodels.Lookup, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookup []dbmodels.Lookup
	err := db.Model(&dbmodels.Lookup{}).Where("lookup_group ~* ?", &lookupstr).Find(&lookup).Error

	if err != nil {
		return nil, constants.ERR_CODE_30, constants.ERR_CODE_30_MSG + " " + err.Error(), err
	}
	return lookup, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
}

// GetPagingLookup ...
func GetPagingLookup(param dto.FilterLookup, offset int, limit int) ([]dbmodels.Lookup, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookup []dbmodels.Lookup
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&lookup).Error
		if err != nil {
			return lookup, 0, err
		}
		return lookup, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQueryCount(db, &total, param, &dbmodels.Lookup{}, "code", errCount)
	// if limit == 0 {
	// 	limit = total
	// }
	go AsyncQuery(db, offset, limit, &lookup, param, "code", errQuery)

	resErrCount := <-errCount
	resErrQuery := <-errQuery

	wg.Done()
	// wg.Done()

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return lookup, 0, resErrCount
	}

	if resErrQuery != nil {
		return lookup, 0, resErrQuery
	}

	return lookup, total, nil
}

// GetLookupFilter ...
func GetLookupFilter(id int) ([]dbmodels.Lookup, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookup []dbmodels.Lookup
	err := db.Model(&dbmodels.Lookup{}).Where("id = ?", &id).First(&lookup).Error

	if err != nil {
		return nil, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG, err
	}
	// else {
	return lookup, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
	// }
}

// GetLookupByGroupName ...
func GetLookupByGroupName(groupName string) ([]dbmodels.Lookup, string, string, error) {

	fmt.Println("Loojkup repository ======>")

	db := GetDbCon()
	db.Debug().LogMode(true)

	fmt.Println("groupName => ", groupName)
	var lookupGroup dbmodels.LookupGroup
	err := db.Model(&dbmodels.LookupGroup{}).Where("name iLike ?", strings.ToUpper(groupName)).Find(&lookupGroup).Error
	if err != nil {
		return nil, constants.ERR_CODE_00, constants.ERR_CODE_51_MSG, err
	}

	var lookup []dbmodels.Lookup
	err = db.Model(&dbmodels.Lookup{}).Where("lookup_group_id = ?", lookupGroup.ID).Find(&lookup).Error

	if err != nil {
		return nil, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG, err
	}
	// else {
	return lookup, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
	// }
}

//SaveLookup ...
func SaveLookup(lookup dbmodels.Lookup) models.NoContentResponse {
	var res models.NoContentResponse
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	db := GetDbCon()
	db.Debug().LogMode(true)

	if lookup.ID < 1 {
		lookup.Code = GenerateLookupCode()
	}

	fmt.Println("Lookup ====> ", lookup)
	lookup.IsViewable = 1
	if r := db.Save(&lookup); r.Error != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}
	return res
}

//GetDistinctLookup ...
func GetDistinctLookup() ([]dbmodels.Lookup, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookup []dbmodels.Lookup
	err := db.Select("DISTINCT lookup_group").Find(&lookup).Error

	// err := db.Model(&dbmodels.Lookup{}).Where("id = ?", &id).First(&lookup).Error

	if err != nil {
		return nil, constants.ERR_CODE_30, constants.ERR_CODE_30_MSG + " " + err.Error(), err
	}
	// else {
	return lookup, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
}

// UpdateLookup ...
func UpdateLookup(updatedlookup dbmodels.Lookup) models.NoContentResponse {
	var res models.NoContentResponse
	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookup dbmodels.Lookup
	err := db.Model(&dbmodels.Lookup{}).Where("id=?", &updatedlookup.ID).First(&lookup).Error
	if err != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res
	}

	lookup.Name = updatedlookup.Name
	lookup.Status = updatedlookup.Status
	lookup.Code = updatedlookup.Code
	lookup.LookupGroupID = updatedlookup.LookupGroupID

	err2 := db.Save(&lookup)
	if err2 != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}

	return res
}

// GenerateLookupCode ...
func GenerateLookupCode() string {
	db := GetDbCon()
	db.Debug().LogMode(true)

	// err := db.Order(order).First(&models)
	var lookup []dbmodels.Lookup
	err := db.Model(&dbmodels.Lookup{}).Order("id desc").First(&lookup).Error
	// err := db.Model(&dbmodels.Brand{}).Where("id = 200000").Order("id desc").First(&brand).Error

	// prefix := loogkupGroup[:2]
	if err != nil {
		return "LK00001"
	}
	if len(lookup) > 0 {
		// fmt.Printf("ini latest code nya : %s \n", brand[0].Code)
		prefix := strings.TrimPrefix(lookup[0].Code, "LK")
		latestCode, err := strconv.Atoi(prefix)
		if err != nil {
			fmt.Printf("error")
			return "LK00001"
		}
		// fmt.Printf("ini latest code nya : %d \n", latestCode)
		wpadding := fmt.Sprintf("%05s", strconv.Itoa(latestCode+1))
		// fmt.Printf("ini pake padding : %s \n", "B"+wpadding)
		return "LK" + wpadding
	}
	return "LK00001"

}

// GetLookupFilter ...
func GetLookupByName(name string) (dbmodels.Lookup, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookup dbmodels.Lookup
	err := db.Model(&dbmodels.Lookup{}).Where("name = ?", &name).First(&lookup).Error

	if err != nil {
		return lookup, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG, err
	}
	// else {
	return lookup, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
	// }
}
