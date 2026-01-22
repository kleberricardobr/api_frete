package interfaces

import "api_frete/models"

type IFreightRepository interface {
	SaveQuote(carriers []models.CarrierInfo) error
	GetCarrierStatistics(limit int) (*models.FreightStatisticsResponse, error)
}
