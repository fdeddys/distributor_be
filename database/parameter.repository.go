package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models"
	"fmt"

)

// GetParameterByNama ...
func GetParameterByNama(nama string) (dbmodels.Parameter, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var parameter dbmodels.Parameter
	err := db.Model(&dbmodels.Parameter{}).Where("name = ?", &nama).First(&parameter).Error

	if err != nil {
		return parameter, constants.ERR_CODE_51, constants.ERR_CODE_51_MSG, err
	}
	// else {
	return parameter, constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, nil
	// }
}

func GetParameter() ([]dbmodels.Parameter, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var parameter []dbmodels.Parameter
	err := db.Model(&dbmodels.Parameter{}).Find(&parameter).Error

	if err != nil {
		return parameter, err
	}
	// else {
	return parameter, nil
	// }
}


// UpdateParam ...
func UpdateApotikParam(updatedParam dbmodels.Parameter) models.NoContentResponse {
	var res models.NoContentResponse
	db := GetDbCon()
	db.Debug().LogMode(true)

	var param dbmodels.Parameter
	err := db.Model(&dbmodels.Parameter{}).Where("id=?", &updatedParam.ID).First(&param).Error
	if err != nil {
		res.ErrCode = "02"
		res.ErrDesc = "Error select data to DB"
		return res
	}

	param.Value = updatedParam.Value
	// param.LastUpdateBy = updatedParam.LastUpdateBy
	// param.LastUpdate = updatedParam.LastUpdate

	err2 := db.Save(&param).Error
	if err2 != nil {
		res.ErrCode = "02"
		res.ErrDesc = "err"
		fmt.Println("Error : " , err2)
		return res
	}

	res.ErrCode = "00"
	res.ErrDesc = "Success"

	return res
}
