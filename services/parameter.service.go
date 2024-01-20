package services

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models"
	// dto "distribution-system-be/models/dto"
	// "time"
	repository "distribution-system-be/database"

)

// ParameterService ...
type ParameterService struct {
}

// GetDataOrderById ...
func (p ParameterService) GetByName(paramName string) dbmodels.Parameter {

	// var res dbmodels.Parameter
	// var err error
	res, errCode, _, _ := database.GetParameterByNama(paramName)
	if errCode == constants.ERR_CODE_00 {
		return res
	}

	return dbmodels.Parameter{ID: 0, Name: "", Value: ""}
}

// GetDataOrderById ...
func (p ParameterService) GetAll() []dbmodels.Parameter {

	// var res dbmodels.Parameter
	// var err error
	res, errCode := database.GetParameter()
	if errCode !=nil {
		return res
	}

	return res
}

// Updateparam ...
func (p ParameterService) UpdateParam(param *dbmodels.Parameter) models.NoContentResponse {
	// param.LastUpdate = time.Now()
	// param.LastUpdateBy = dto.CurrUser
	
	res := repository.UpdateApotikParam(*param)

	return res
}