package mocks

import (
	"api_frete/models"
)

type MockFreightRepository struct {
	SaveQuoteFunc            func(carriers []models.CarrierInfo) error
	GetCarrierStatisticsFunc func(limit int) (*models.FreightStatisticsResponse, error)
}

func (m *MockFreightRepository) SaveQuote(carriers []models.CarrierInfo) error {
	if m.SaveQuoteFunc != nil {
		return m.SaveQuoteFunc(carriers)
	}
	return nil
}

func (m *MockFreightRepository) GetCarrierStatistics(limit int) (*models.FreightStatisticsResponse, error) {
	if m.GetCarrierStatisticsFunc != nil {
		return m.GetCarrierStatisticsFunc(limit)
	}
	return &models.FreightStatisticsResponse{}, nil
}
