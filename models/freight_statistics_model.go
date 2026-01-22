package models

type FreightStatisticsResponse struct {
	Carrier  []CarrierStatistics `json:"carrier"`
	MinPrice float64             `json:"min_price"`
	MaxPrice float64             `json:"max_price"`
}

type CarrierStatistics struct {
	Name           string  `json:"name"`
	TotalFreight   float64 `json:"total_freight"`
	AverageFreight float64 `json:"average_freight"`
	QtyResults     int     `json:"qty_results"`
}
