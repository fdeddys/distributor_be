package dtovendors

// SalesOrderResultDto ...
type SalesOrderResultDto struct {
	Success     SalesOrderSuccesstDto `json:"success"`
	AccessToken string                `json:"access_token"`
}

// SalesOrderSuccesstDto ...
type SalesOrderSuccesstDto struct {
	Page         int             `json:"page"`
	PerPage      int             `json:"per_page"`
	TotalPages   int             `json:"total_pages"`
	TotalRecords int             `json:"total_records"`
	Data         []SalesOrderDto `json:"data"`
}

// SalesOrderDto ...
type SalesOrderDto struct {
	Code                                    string           `json:"code"`
	CustomerCode                            string           `json:"customer_code"`
	WarehouseCode                           string           `json:"warehouse_code"`
	State                                   string           `json:"state"`
	Currency                                string           `json:"currency"`
	ProjectCode                             string           `json:"project_code"`
	RoundingAmount                          string           `json:"rounding_amount"`
	SalesEmployeeCode                       string           `json:"sales_employee_code"`
	ShipFromAddressContactPersonName        string           `json:"ship_from_address_contact_person_name"`
	ShipFromAddressContactPersonPhoneNumber string           `json:"ship_from_address_contact_person_phone_number"`
	ShipFromAddressLine1                    string           `json:"ship_from_address_line1"`
	ShipFromAddressLine2                    string           `json:"ship_from_address_line2"`
	ShipFromAddressLine3                    string           `json:"ship_from_address_line3"`
	ShipFromAddressLine4                    string           `json:"ship_from_address_line4"`
	ShipFromAddressCountryCode              string           `json:"ship_from_address_country_code"`
	ShipFromAddressPostcode                 string           `json:"ship_from_address_postcode"`
	ShipToAddressContactPersonName          string           `json:"ship_to_address_contact_person_name"`
	ShipToAddressContactPersonPhoneNumber   string           `json:"ship_to_address_contact_person_phone_number"`
	ShipToAddressLine1                      string           `json:"ship_to_address_line1"`
	ShipToAddressLine2                      string           `json:"ship_to_address_line2"`
	ShipToAddressLine3                      string           `json:"ship_to_address_line3"`
	ShipToAddressLine4                      string           `json:"ship_to_address_line4"`
	ShipToAddressCountryCode                string           `json:"ship_to_address_country_code"`
	ShipToAddressPostcode                   string           `json:"ship_to_address_postcode"`
	BillToAddressContactPersonName          string           `json:"bill_to_address_contact_person_name"`
	BillToAddressContactPersonPhoneNumber   string           `json:"bill_to_address_contact_person_phone_number"`
	BillToAddressLine1                      string           `json:"bill_to_address_line1"`
	BillToAddressLine2                      string           `json:"bill_to_address_line2"`
	BillToAddressLine3                      string           `json:"bill_to_address_line3"`
	BillToAddressline4                      string           `json:"bill_to_address_line4"`
	BillToAddressCountryCode                string           `json:"bill_to_address_country_code"`
	BillToAddressPostcode                   string           `json:"bill_to_address_postcode"`
	CompanyAddressLine1                     string           `json:"company_address_line1"`
	CompanyAddressLine2                     string           `json:"company_address_line2"`
	CompanyAddressLine3                     string           `json:"company_address_line3"`
	CompanyAddressLine4                     string           `json:"company_address_line4"`
	CompanyAddressCountryCode               string           `json:"company_address_country_code"`
	CompanyAddressPostcode                  string           `json:"company_address_postcode"`
	TransactionAt                           string           `json:"transaction_at"`
	EstimatedDeliveryAt                     string           `json:"estimated_delivery_at"`
	CreatedAt                               string           `json:"created_at"`
	UpdatedAt                               string           `json:"updated_at"`
	SalesOrderLines                         []SalesOrderLine `json:"sales_order_lines"`
}

// SalesOrderLine ...
type SalesOrderLine struct {
	ItemCode          string `json:"item_code"`
	Quantity          string `json:"quantity"`
	Price             string `json:"price"`
	UOMCode           string `json:"unit_of_measurement_code"`
	Description       string `json:"description"`
	ProjectCode       string `json:"project_code"`
	SalesEmployeeCode string `json:"sales_employee_code"`
	DiscountCode      string `json:"discount_code"`
	TaxCode           string `json:"tax_code"`
}

// ParamGetSoDto ...
type ParamGetSoDto struct {
	AccessToken     string
	Page            string
	PerPage         string
	SalesOrderCodes string
}
