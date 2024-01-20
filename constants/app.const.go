package constants

const (
	TokenSecretKey        = "OasI$_sEcrET_key$"
	TokenExpiredInMinutes = 8 * 60 * 60
	VERSION               = "B180622"
	VERSION_DATE          = "28OKT2021"
)

// ERR code Global
const (
	ERR_CODE_00     = "00"
	ERR_CODE_00_MSG = "SUCCESS.."

	ERR_CODE_03     = "03"
	ERR_CODE_03_MSG = "Error, unmarshall body Request"

	ERR_CODE_05     = "05"
	ERR_CODE_05_MSG = "Error,Cannot get parameter query"

	ERR_CODE_06     = "06"
	ERR_CODE_06_MSG = "Error,Cannot get File"

	ERR_CODE_07     = "07"
	ERR_CODE_07_MSG = "Error,Cannot parse File"
)

// ERR code Global
const (
	ERR_CODE_30     = "30"
	ERR_CODE_30_MSG = "Failed save data to DB"
)

const (
	ERR_CODE_50     = "50"
	ERR_CODE_50_MSG = "Invalid username / password"

	ERR_CODE_51     = "51"
	ERR_CODE_51_MSG = "Error connection to database"

	ERR_CODE_53     = "53"
	ERR_CODE_53_MSG = "Failed generate token !"

	ERR_CODE_54     = "54"
	ERR_CODE_54_MSG = "Invalid Authorization !"

	ERR_CODE_55     = "55"
	ERR_CODE_55_MSG = "Token expired !"

	ERR_CODE_61     = "61"
	ERR_CODE_61_MSG = "User not found !"

	ERR_CODE_62     = "62"
	ERR_CODE_62_MSG = "Password not match !"

	ERR_CODE_63     = "63"
	ERR_CODE_63_MSG = "Failed Update password !"
)

//nDeskKey ...
const (
	DesKey = "abcdefghijklmnopqrstuvwxyz012345"
)

// ERROR FROM VENDOR
const (
	ERR_CODE_70     = "70"
	ERR_CODE_70_MSG = "ERROR FROM VENDOR "
)

const (
	ERR_CODE_80     = "80"
	ERR_CODE_80_MSG = "Failed save to database"

	ERR_CODE_81     = "81"
	ERR_CODE_81_MSG = "Failed get data from database"
)

const (
	ERR_CODE_90     = "90"
	ERR_CODE_90_MSG = "Failed get from Body Request"

	ERR_CODE_95     = "95"
	ERR_CODE_95_MSG = "NO Order found for payment !"

	ERR_CODE_96     = "96"
	ERR_CODE_96_MSG = "Payment total order / return / detail not MATCH  !"

	ERR_CODE_99     = "99"
	ERR_CODE_99_MSG = "Error business process"
)

const (
	ERR_CODE_40     = "40"
	ERR_CODE_40_MSG = "Data not found"

	ERR_CODE_41     = "41"
	ERR_CODE_41_MSG = "STATUS NOT VALID, PLEASE REFRESH YOUR DATA !"
)

// STATUS Sales Order
// 10 = new order
// 20 = approve
// 30 = reject
// 40 = INVOICE
// 50 = PAID
// 60 = Reject Payment

const (
	STATUS_APPROVE        = 20
	STATUS_NEW            = 10
	STATUS_REJECT         = 30
	STATUS_PAID           = 50
	STATUS_REJECT_PAYMENT = 60
)

const (
	PARAMETER_TAX_VALUE     = "tax"
	PARAMETER_APOTEKER_NAME = "apoteker"
	PARAMETER_SIA           = "SIA"
	PARAMETER_SIPA          = "SIPA"
)

const (
	HEADER_PAYMENT_CASH     = "PC"
	HEADER_PAYMENT_CREDIT   = "PY"
	HEADER_PAYMENT_SUPPLIER = "PS"
)
