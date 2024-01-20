package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"fmt"
)

// AddSequence ...
func AddSequence(month, year, header string) (errcd string, newNumb int32, errdesc string) {
	db := GetDbCon()

	var seq dbmodels.Sequence
	var urut int32

	db.Where("year = ? and month = ? and subj = ?", year, month, header).First(&seq)

	fmt.Println("Seq ID == ", seq.ID)
	if seq.ID == 0 {
		errCode, errDesc := NewSequence(year, month, header)
		return errCode, 1, errDesc
	}
	seq.Seq++
	urut = seq.Seq
	db.Save(&seq)
	// db.Model(&seq).Where("id = ?", seq.ID).Update("seq", urut)

	// var code = ""
	// code = constants.ERR_CODE_00
	return constants.ERR_CODE_00, urut, constants.ERR_CODE_00_MSG
}

// NewSequence ...
func NewSequence(year, month, header string) (errcode string, errdesc string) {
	db := GetDbCon()

	var seq dbmodels.Sequence
	seq.Month = month
	seq.Subject = header
	seq.Year = year
	seq.Seq = 1
	err := db.Save(&seq)
	if err.Error != nil {
		return constants.ERR_CODE_80, constants.ERR_CODE_80_MSG
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
