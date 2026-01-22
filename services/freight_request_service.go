package services

import (
	"api_frete/interfaces"
	"api_frete/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const URL_FRETE = "https://sp.freterapido.com/api/v3/quote/simulate"

type FreightService struct {
	Repository interfaces.IFreightRepository
	ClientHttp interfaces.IHTTPClient
}

func NewFreightService(repo interfaces.IFreightRepository,
	client interfaces.IHTTPClient) interfaces.IFreightService {
	return &FreightService{
		Repository: repo,
		ClientHttp: client,
	}
}

func (s *FreightService) SendFreightQuote(data *models.FreightRequest) (error, *models.FreightCarrierResponse) {
	request, err := s.createFreightQuoteRequest(data)

	if err != nil {
		return fmt.Errorf("erro on prepare data to request: %v", err), nil
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error on JSON convertion: %v", err), nil
	}

	req, err := http.NewRequest("POST", URL_FRETE, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error on create request: %v", err), nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.ClientHttp.Do(req)
	if err != nil {
		return fmt.Errorf("erro on sending request: %v", err), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro on response: status %d", resp.StatusCode), nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error on body reading: %v", err), nil
	}

	var response models.FreightQuoteResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return fmt.Errorf("error on response parsing: %v", err), nil
	}

	carriers := s.convertToCarrierResponse(response)
	err = s.Repository.SaveQuote(carriers.Carrier)
	if err != nil {
		return fmt.Errorf("error on save to database: %v", err), nil
	}

	return nil, carriers
}

func (s *FreightService) createFreightQuoteRequest(data *models.FreightRequest) (*models.FreightQuoteRequest, error) {
	zipCode, err := strconv.Atoi(data.Recipient.Address.Zipcode)
	if err != nil {
		return nil, fmt.Errorf("Failed on get recipient zipcode: %v", err)
	}

	volumes := []models.VolumeQuote{}

	for _, v := range data.Volumes {

		if v.Amount <= 0 {
			return nil, fmt.Errorf("Amount value must be bigger then zero")
		}

		volume := models.VolumeQuote{
			Category:      strconv.Itoa(v.Category),
			Amount:        v.Amount,
			UnitaryWeight: v.UnitaryWeight,
			Price:         v.Price,
			SKU:           v.SKU,
			Height:        v.Height,
			Width:         v.Width,
			Length:        v.Length,
			UnitaryPrice:  v.Price / float64(v.Amount),
		}

		volumes = append(volumes, volume)
	}

	return &models.FreightQuoteRequest{
		Shipper: models.Shipper{
			RegisteredNumber: os.Getenv("REG_NUMBER"),
			Token:            os.Getenv("TOKEN"),
			PlatformCode:     os.Getenv("SYS_CODE"),
		},
		Recipient: models.RecipientZip{
			Zipcode: zipCode,
		},
		Dispatchers: []models.Dispatcher{
			{
				RegisteredNumber: os.Getenv("REG_NUMBER"),
				Zipcode:          s.getDispZipcode(),
				TotalPrice:       s.calculateTotalPrice(volumes),
				Volumes:          volumes,
			},
		},
		SimulationType: []int{0},
	}, nil
}

func (s *FreightService) convertToCarrierResponse(response models.FreightQuoteResponse) *models.FreightCarrierResponse {
	carriers := []models.CarrierInfo{}

	for _, dispatcher := range response.Dispatchers {
		for _, offer := range dispatcher.Offers {
			carrier := models.CarrierInfo{
				Name:     offer.Carrier.Name,
				Service:  offer.Service,
				Deadline: offer.DeliveryTime.Days,
				Price:    offer.FinalPrice,
			}
			carriers = append(carriers, carrier)
		}
	}

	return &models.FreightCarrierResponse{
		Carrier: carriers,
	}
}

func (s *FreightService) getDispZipcode() int {
	var zipcode int
	fmt.Sscanf(os.Getenv("DISP_ZIPCODE"), "%d", &zipcode)
	return zipcode
}

func (s *FreightService) calculateTotalPrice(volumes []models.VolumeQuote) float64 {
	total := 0.0
	for _, v := range volumes {
		total += v.Price
	}
	return total
}
