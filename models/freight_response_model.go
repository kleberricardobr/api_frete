package models

type FreightQuoteResponse struct {
	Dispatchers []DispatcherResponse `json:"dispatchers"`
}

type DispatcherResponse struct {
	Offers []OfferResponse `json:"offers"`
}

type OfferResponse struct {
	Carrier      CarrierResponse      `json:"carrier"`
	Service      string               `json:"service"`
	FinalPrice   float64              `json:"final_price"`
	DeliveryTime DeliveryTimeResponse `json:"delivery_time"`
}

type CarrierResponse struct {
	Name string `json:"name"`
}

type DeliveryTimeResponse struct {
	Days int `json:"days"`
}
