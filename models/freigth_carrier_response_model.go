package models

type FreightCarrierResponse struct {
	Carrier []CarrierInfo `json:"carrier"`
}

type CarrierInfo struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}
