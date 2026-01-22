package services

import (
	"api_frete/interfaces"
	"api_frete/models"
)

type MetricsService struct {
	Repository interfaces.IFreightRepository
}

func NewMetricsService(repo interfaces.IFreightRepository) interfaces.IFreightMetricsInterface {
	return &MetricsService{
		Repository: repo,
	}
}

func (s *MetricsService) GetMetrics(lastQuotes int) (*models.FreightStatisticsResponse, error) {
	stats, err := s.Repository.GetCarrierStatistics(lastQuotes)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
