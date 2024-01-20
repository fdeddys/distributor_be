package dto

// DBoardQtyDto ...
type DBoardQtyDto struct {
	OrderCount     int `json:"order_count"`
	InternalStatus int `json:"internal_status"`
}

// DBoardQtyOrderDto ...
type DBoardQtyOrderDto struct {
	NewOrder int `json:"neworder"`
	Order    int `json:"order"`
	Payment  int `json:"payment"`
	Complete int `json:"complete"`
	Reject   int `json:"reject"`
}
