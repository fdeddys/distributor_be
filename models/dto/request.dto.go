package dto

// RequestHeader ...
type RequestHeader struct {
	Signature     string `valid:"Required";json:"signature"`
	Timestamp     string `valid:"Required";json:"timestamp"`
	Authorization string `valid:"Required";json:"authorization"`
}
