package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
)

//GetLookupGroup ...
func GetLookupGroup() ([]dbmodels.LookupGroup, string, string, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var lookupGroup []dbmodels.LookupGroup
	err := db.Model(&dbmodels.Lookup{}).Find(&lookupGroup).Error

	if err != nil {
		return nil, constants.ERR_CODE_30, constants.ERR_CODE_30_MSG + " " + err.Error(), err
	}
	return lookupGroup, "00", "success", nil
}
