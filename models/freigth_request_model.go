package models

type FreightQuoteRequest struct {
	Shipper        Shipper      `json:"shipper"`
	Recipient      RecipientZip `json:"recipient"`
	Dispatchers    []Dispatcher `json:"dispatchers"`
	SimulationType []int        `json:"simulation_type"`
}

type Shipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type RecipientZip struct {
	Zipcode int `json:"zipcode"`
}

type Dispatcher struct {
	RegisteredNumber string        `json:"registered_number"`
	Zipcode          int           `json:"zipcode"`
	TotalPrice       float64       `json:"total_price"`
	Volumes          []VolumeQuote `json:"volumes"`
}

type VolumeQuote struct {
	Category      string  `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Price         float64 `json:"price"`
	SKU           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
}
