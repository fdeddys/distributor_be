package models

// "time"

// ResponsePagination .. pagination
type ResponsePagination struct {
	TotalRow int         `json:"totalRow"`
	Page     int         `json:"page"`
	Count    int         `json:"count"`
	Contents interface{} `json:"contents"`
	Error    string      `json:"error"`
}

// ContentResponse ...
type ContentResponse struct {
	ErrCode  string      `json:"errCode"`
	ErrDesc  string      `json:"errDesc"`
	Contents interface{} `json:"contents"`
}

// NoContentResponse ...
type NoContentResponse struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
}

// Response ...
type Response struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
}

type ResponseCheckPrice struct {
	ErrCode string  `json:"errCode"`
	ErrDesc string  `json:"errDesc"`
	Price   int64   `json:"price"`
	Disc1   int64   `json:"disc1"`
	Disc2   int64   `json:"disc2"`
	Hpp     float32 `json:"hpp"`
}

// Merchant ...
type ResponseMerchant struct {
	Data interface{} `json:"data"`
	// ID           				int8  	`json:"id"`
	// Code     					string 	`json:"code"`
	// Name        				string 	`json:"name"`
	// IssuerCode     				string 	`json:"issuerCode"`
	// Top       					int    	`json:"top"`
	// Status 						int 	`json:"status"`
	// LastUpdateBy				string	`json:"last_update_by"`
	// LastUpdate	time.Time				`json:"last_update"`
	// IssuerName					string  `json:"issuerName"`
	// Data         interface{} `json:"data"`
	// ID           int8        `json:"id"`
	// Code         string      `json:"code"`
	// Name         string      `json:"name"`
	// IssuerCode   string      `json:"issuerCode"`
	// Top          int         `json:"top"`
	// Status       int         `json:"status"`
	// LastUpdateBy string      `json:"last_update_by"`
	// LastUpdate   time.Time   `json:"last_update"`
	// IssuerName   string      `json:"issuerName"`
}

//
type ResponseIssuer struct {
	Data interface{} `json:"data"`
}

// Response Supplier
// type ResponseSupplier struct {
// 	Data interface{} `json:"data"`
// }

// ResponseUser ...
type ResponseUser struct {
	UserName     string `json:"userName" gorm:"column:user_name"`
	Email        string `json:"email" gorm:"column:email"`
	Status       int    `json:"status" gorm:"column:status"`
	SupplierCode string `json:"supplierCode" gorm:"column:supplier_code"`
}

// ResponseWarehouseBySupplierId
type ResponseSupplierWarehouse struct {
	ErrCode string      `json:"errCode"`
	ErrDesc string      `json:"errDesc"`
	Data    interface{} `json:"data"`
}

type ResponseSupplier struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Code    string `json:"code"`
}

type ResponseSupplierGroup struct {
	Data interface{} `json:"data"`
}

type SupplierListResponse struct {
	Supplier []Supplier `json:"suppliers"`
}

type Supplier struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	TechRequirement []string `json:"tech_requirement"`
	DocRequirement  []string `json:"doc_requirement"`
	MinOrder        int      `json:"min_order"`
	DownPayment     string   `json:"down_payment"`
	TypeNoo         bool     `json:"type_noo"`
}

// type Logo struct {
// 	Url   string `json:"url"`
// 	Small Url    `json:"small"`
// 	Thumb Url    `json:"thumb"`
// 	Big   Url    `json:"big"`
// }

// type Url struct {
// 	Url string `json:"url"`
// }

type ResponseSFA struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseCheckOrder struct {
	ErrCode string `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Data    string `json:"data"`
}

// type ResponseCheckOrder struct {
// 	ErrCode 		string	`json:"errCode"`
// 	ErrDesc 		string  `json:"errDesc"`
// 	Data			string	`json:"data"`
// }

type ResponseUpload struct {
	ErrCode  string `json:"errCode"`
	ErrDesc  string `json:"errDesc"`
	FileName string `json:"fileName"`
}

type ResponseSupplierProduct struct {
	ID           int64  `json:"id"`
	Price        int    `json:"price"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

// RequestUpdateStatus ...
type RequestUpdateStatus struct {
	SoNumber   string `json:"soNumber"`
	StatusDesc string `json:"statusDesc"`
}

type ResponseMerchantBySupplier struct {
	Message         string   `json:"message"`
	Description     string   `json:"description"`
	SupplierID      int64    `json:"supplier_id"`
	TechRequirement []string `json:"tech_requirement"`
	DocRequirement  []string `json:"doc_requirement"`
}

type ResponseFollowUpOrder struct {
	ErrCode string      `json:"errCode"`
	ErrDesc string      `json:"errDesc"`
	Total   int         `json:"total"`
	Data    interface{} `json:"data"`
}

type ResponseStatus struct {
	ErrCode        string `json:"errCode"`
	ErrDesc        string `json:"errDesc"`
	ApprovalStatus bool   `json:"approvalStatus"`
}

type ResponseReceiveCheckPrice struct {
	Price int64 `json:"price"`
}
