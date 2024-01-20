package dtovendors

//MerchantDTO ..
type MerchantDTO struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	StateAction string `json:"state_action"`
	Code        string `json:"code"`
}

// MerchantResponse ...
type MerchantResponse struct {
	Success     SuccessResponse `json:"success"`
	Error       ErrorResponse   `json:"error"`
	AccessToken string          `json:"access_token"`
}

//SuccessResponse ..
type SuccessResponse struct {
	WarehouseCode string `json:"warehouse_code"`
}

//ErrorResponse ..
type ErrorResponse struct {
	Message string `json:"message"`
}


// Response Token Uki 
type ResponTokenUki struct {
	AccessToken string `json:"access_token"`
}


// CustomerResponse ...
type CustomerResponse struct {
	Success     SuccessCustomerResponse `json:"success"`
	Error       ErrorResponse   `json:"error"`
	AccessToken string          `json:"access_token"`
}

//SuccessResponse ..
type SuccessCustomerResponse struct {
	CustomerCode string `json:"customer_code"`
}