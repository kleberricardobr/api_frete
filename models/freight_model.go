package models

type FreightRequest struct {
	Recipient Recipient `json:"recipient"`
	Volumes   []Volume  `json:"volumes"`
}

type Recipient struct {
	Address Address `json:"address"`
}

type Address struct {
	Zipcode string `json:"zipcode"`
}

type Volume struct {
	Category      int     `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Price         float64 `json:"price"`
	SKU           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}
