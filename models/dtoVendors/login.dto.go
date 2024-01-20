package dtovendors

// LoginRequest ..
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponseVendor ..
type LoginResponseVendor struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

// LoginResponse ..
type LoginResponse struct {
	Token   string `json:"token"`
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
}
