package interfaces

import "api_frete/models"

type IFreightService interface {
	SendFreightQuote(data *models.FreightRequest) (error, *models.FreightCarrierResponse)
}
