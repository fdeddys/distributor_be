package services

import (
	"distribution-system-be/constants"
	repository "distribution-system-be/database"
	"distribution-system-be/models"
)

// LookupGroupService ...
type LookupGroupService struct {
}

// GetLookupGroup ...
func (h LookupGroupService) GetLookupGroup() models.ContentResponse {
	var res models.ContentResponse

	data, errCode, errDesc, err := repository.GetLookupGroup()
	res.ErrCode = errCode
	res.ErrDesc = errDesc
	res.Contents = data

	if err != nil {
		res.Contents = nil
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
	}
	return res
}
