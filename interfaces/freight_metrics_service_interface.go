package interfaces

import "api_frete/models"

type IFreightMetricsInterface interface {
	GetMetrics(lastQuotes int) (*models.FreightStatisticsResponse, error)
}
