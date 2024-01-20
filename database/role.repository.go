package database

import (
	constants "distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"fmt"
	"log"
	"sync"

	dto "distribution-system-be/models/dto"

	"github.com/jinzhu/gorm"
)

// type TotalRows struct {
// 	Total int `gorm:"column(count)"`
// }

// GetRole ...
func GetRole(param dto.FilterName, offset int, limit int) ([]dbmodels.Role, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var role []dbmodels.Role
	var total int
	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&role).Error
		if err != nil {
			return role, 0, err
		}
		return role, 0, nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	// go asyncQueryRoles(db, offset, limit, &role, errQuery)
	// go asyncQueryCountRoles(db, &total, errCount)
	go AsyncQuery(db, offset, limit, &role, param, "name", errQuery)
	go AsyncQueryCount(db, &total, param, &dbmodels.Role{}, "name", errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return role, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return role, 0, resErrCount
	}
	return role, total, nil
}

// asyncQueryRoles ...
func asyncQueryRoles(db *gorm.DB, offset int, limit int, user *[]dbmodels.Role, resChan chan error) {
	err := db.Offset(offset).Limit(limit).Find(&user).Error
	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// asyncQueryCountRoles ...
func asyncQueryCountRoles(db *gorm.DB, total *int, resChan chan error) {
	err := db.Model(&dbmodels.User{}).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// UpdateRole ...
func UpdateRole(updatedRole dbmodels.Role) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	var role dbmodels.Role
	err := db.Model(&dbmodels.Role{}).Where("id=?", &updatedRole.ID).First(&role).Error
	if err != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res
	}

	role.Name = updatedRole.Name
	role.Description = updatedRole.Description

	err2 := db.Save(&role)
	if err2 != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res
	}

	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG

	return res
}

//SaveRole ...
func SaveRole(role dbmodels.Role) (models.NoContentResponse, int64) {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	// brand.Code = GenerateBrandCode()
	r := db.Save(&role)

	fmt.Println("Last id sav ==> ", role.ID)
	if r.Error != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res, role.ID
	}

	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG
	return res, role.ID
}
