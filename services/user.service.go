package services

import (
	repository "distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"fmt"

	// router "distribution-system-be/routers"
	"time"

	kons "distribution-system-be/constants"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserService ...
type UserService struct {
}

// GetUserFilterPaging ...
func (h UserService) GetUserFilterPaging(param dto.FilterUser, page int, limit int) models.ResponsePagination {

	var res models.ResponsePagination

	offset := (page - 1) * limit
	data, totalData, err := repository.GetUserTransaction(param, offset, limit)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.Contents = data
	res.TotalRow = totalData
	res.Page = page
	res.Count = len(data)

	return res
}

// AuthLogin ...
func (h UserService) AuthLogin(userDto *dto.LoginRequestDto) dto.LoginResponseDto {
	var res dto.LoginResponseDto
	if userDto.Username == "" {
		res.ErrCode = kons.ERR_CODE_50
		res.ErrDesc = kons.ERR_CODE_50_MSG
		return res
	}

	if userDto.Password == "" {
		res.ErrCode = kons.ERR_CODE_50
		res.ErrDesc = kons.ERR_CODE_50_MSG
		return res
	}

	user, err := repository.GetUserByName(userDto.Username)
	if err != nil {
		res.ErrCode = kons.ERR_CODE_51
		res.ErrDesc = kons.ERR_CODE_51_MSG
		return res
	}
	// fmt.Println("USER ---> ", user)
	if user.ID == 0 {
		res.ErrCode = kons.ERR_CODE_50
		res.ErrDesc = kons.ERR_CODE_50_MSG
		return res
	}

	if user.Password != userDto.Password {
		res.ErrCode = kons.ERR_CODE_50
		res.ErrDesc = kons.ERR_CODE_50_MSG
		return res
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["user"] = user.UserName
	claims["userId"] = fmt.Sprintf("%v", (user.ID))
	// now := time.Now()
	// claims["logTm"] = now
	// claims["supplierCode"] = user.SupplierCode

	unixNano := time.Now().UnixNano()
	umillisec := unixNano / 1000000
	timeToString := fmt.Sprintf("%v", umillisec)
	fmt.Println("token Created ", timeToString)
	claims["tokenCreated"] = timeToString

	token, err := sign.SignedString([]byte(kons.TokenSecretKey))
	if err != nil {
		res.ErrCode = kons.ERR_CODE_53
		res.ErrDesc = kons.ERR_CODE_53_MSG
		res.Token = ""
		return res
	}

	res.ErrCode = kons.ERR_CODE_00
	res.ErrDesc = kons.ERR_CODE_00_MSG
	res.Token = token

	return res
}

// ChangePass ...
func (h UserService) ChangePass(userDto *dto.ChangePassRequestDto) models.Response {

	var res models.Response

	var currUser string
	currUser = dto.CurrUser

	fmt.Println("Cur user -> ", currUser)

	user, err := repository.GetByUsername(currUser)
	if err != nil {
		res.ErrCode = kons.ERR_CODE_51
		res.ErrDesc = fmt.Sprintf("%s [%s]", kons.ERR_CODE_51_MSG, err.Error())
		return res
	}

	if user.ID == 0 {
		res.ErrCode = kons.ERR_CODE_61
		res.ErrDesc = kons.ERR_CODE_61_MSG
		return res
	}

	if user.Password != userDto.OldPassword {
		res.ErrCode = kons.ERR_CODE_62
		res.ErrDesc = kons.ERR_CODE_62_MSG
		return res
	}

	if _, err2 := repository.UpdatePassword(currUser, userDto.NewPassword); err2 != nil {
		res.ErrCode = kons.ERR_CODE_63
		res.ErrDesc = kons.ERR_CODE_63_MSG
		return res
	}

	res.ErrCode = kons.ERR_CODE_00
	res.ErrDesc = kons.ERR_CODE_00_MSG

	return res
}

func (h UserService) SaveUser(user *dbmodels.User) models.ContentResponse {
	user.LastUpdate = time.Now()
	user.LastUpdateBy = dto.CurrUser

	res := repository.SaveUser(*user)
	return res
}

// UpdateUser ...
func (h UserService) UpdateUser(user *dbmodels.User) models.NoContentResponse {
	user.LastUpdate = time.Now()
	user.LastUpdateBy = dto.CurrUser
	// var updatedBrand dbmodels.Brand
	// updatedBrand.ID = brand.ID
	// updatedBrand.Name = brand.Name
	// updatedBrand.Status = brand.Status
	// updatedBrand.LastUpdateBy = brand.LastUpdateBy
	// updatedBrand.LastUpdate = brand.LastUpdate
	// updatedBrand.Code = brand.Code

	res := repository.UpdateUser(*user)

	return res
}

// ResetUser ...
func (h UserService) ResetUser(idUser int64) models.ContentResponse {

	return repository.ResetPassword(idUser)
}
