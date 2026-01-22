package services

import (
	"api_frete/mocks"
	"api_frete/models"
	"errors"
	"net/http"
	"os"
	"testing"
)

func setupEnvVars() {
	os.Setenv("REG_NUMBER", "25438296000158")
	os.Setenv("TOKEN", "test-token")
	os.Setenv("SYS_CODE", "test-code")
	os.Setenv("DISP_ZIPCODE", "29161376")
}

func cleanupEnvVars() {
	os.Unsetenv("REG_NUMBER")
	os.Unsetenv("TOKEN")
	os.Unsetenv("SYS_CODE")
	os.Unsetenv("DISP_ZIPCODE")
}

func TestFreightService_SendFreightQuote_Success(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{
		SaveQuoteFunc: func(carriers []models.CarrierInfo) error {
			if len(carriers) != 2 {
				t.Errorf("Expected 2 carriers, got %d", len(carriers))
			}
			return nil
		},
	}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{
				Category:      7,
				Amount:        1,
				UnitaryWeight: 5,
				Price:         349,
				SKU:           "abc-teste-123",
				Height:        0.2,
				Width:         0.2,
				Length:        0.2,
			},
		},
	}

	mockResponse := `{
        "dispatchers": [{
            "offers": [
                {
                    "carrier": {"name": "CORREIOS"},
                    "service": "PAC",
                    "delivery_time": {"days": 5},
                    "final_price": 25.50
                },
                {
                    "carrier": {"name": "JADLOG"},
                    "service": "Expresso",
                    "delivery_time": {"days": 3},
                    "final_price": 35.00
                }
            ]
        }]
    }`

	mockHTTPClient := mocks.NewMockHTTPClientWithResponse(200, mockResponse)

	service := NewFreightService(mockRepo, mockHTTPClient).(*FreightService)

	err, response := service.SendFreightQuote(request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(response.Carrier) != 2 {
		t.Errorf("Expected 2 carriers in response, got %d", len(response.Carrier))
	}

	if response.Carrier[0].Name != "CORREIOS" {
		t.Errorf("Expected first carrier 'CORREIOS', got '%s'", response.Carrier[0].Name)
	}

	if response.Carrier[0].Price != 25.50 {
		t.Errorf("Expected first carrier price 25.50, got %.2f", response.Carrier[0].Price)
	}

	if response.Carrier[1].Deadline != 3 {
		t.Errorf("Expected second carrier deadline 3 days, got %d", response.Carrier[1].Deadline)
	}
}

func TestFreightService_SendFreightQuote_HTTPError(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{
				Category: 7,
				Amount:   1,
				Price:    100,
			},
		},
	}

	mockHTTPClient := &mocks.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	service := NewFreightService(mockRepo, mockHTTPClient).(*FreightService)

	err, response := service.SendFreightQuote(request)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if response != nil {
		t.Error("Expected nil response, got data")
	}

	if !containsString(err.Error(), "network error") {
		t.Errorf("Expected error message to contain 'network error', got '%s'", err.Error())
	}
}

func TestFreightService_SendFreightQuote_BadStatusCode(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{
				Category: 7,
				Amount:   1,
				Price:    100,
			},
		},
	}

	mockHTTPClient := mocks.NewMockHTTPClientWithResponse(500, `{"error": "Internal Server Error"}`)

	service := NewFreightService(mockRepo, mockHTTPClient).(*FreightService)

	err, response := service.SendFreightQuote(request)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if response != nil {
		t.Error("Expected nil response, got data")
	}

	if !containsString(err.Error(), "status 500") {
		t.Errorf("Expected error message to contain 'status 500', got '%s'", err.Error())
	}
}

func TestFreightService_SendFreightQuote_InvalidJSON(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{
				Category: 7,
				Amount:   1,
				Price:    100,
			},
		},
	}

	mockHTTPClient := mocks.NewMockHTTPClientWithResponse(200, `{invalid json}`)

	service := NewFreightService(mockRepo, mockHTTPClient).(*FreightService)

	err, response := service.SendFreightQuote(request)

	if err == nil {
		t.Error("Expected parsing error, got nil")
	}

	if response != nil {
		t.Error("Expected nil response, got data")
	}

	if !containsString(err.Error(), "parsing") {
		t.Errorf("Expected error message to contain 'parsing', got '%s'", err.Error())
	}
}

func TestFreightService_SendFreightQuote_RepositoryError(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{
		SaveQuoteFunc: func(carriers []models.CarrierInfo) error {
			return errors.New("database connection failed")
		},
	}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{
				Category: 7,
				Amount:   1,
				Price:    100,
			},
		},
	}

	mockResponse := `{
        "dispatchers": [{
            "offers": [{
                "carrier": {"name": "CORREIOS"},
                "service": "PAC",
                "delivery_time": {"days": 5},
                "final_price": 25.50
            }]
        }]
    }`

	mockHTTPClient := mocks.NewMockHTTPClientWithResponse(200, mockResponse)

	service := NewFreightService(mockRepo, mockHTTPClient).(*FreightService)

	err, response := service.SendFreightQuote(request)

	if err == nil {
		t.Error("Expected database error, got nil")
	}

	if response != nil {
		t.Error("Expected nil response, got data")
	}

	if !containsString(err.Error(), "database") {
		t.Errorf("Expected error message to contain 'database', got '%s'", err.Error())
	}
}

func TestFreightService_CreateFreightQuoteRequest_Success(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{
				Category:      7,
				Amount:        2,
				UnitaryWeight: 5,
				Price:         100,
				SKU:           "test-123",
				Height:        0.2,
				Width:         0.3,
				Length:        0.4,
			},
		},
	}

	mockResponse := `{
        "dispatchers": [{
            "offers": [{
                "carrier": {"name": "CORREIOS"},
                "service": "PAC",
                "delivery_time": {"days": 5},
                "final_price": 25.50
            }]
        }]
    }`

	mockHTTPClient := mocks.NewMockHTTPClientWithResponse(200, mockResponse)

	service := NewFreightService(mockRepo, mockHTTPClient).(*FreightService)

	quoteRequest, err := service.createFreightQuoteRequest(request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if quoteRequest.Recipient.Zipcode != 1311000 {
		t.Errorf("Expected zipcode 1311000, got %d", quoteRequest.Recipient.Zipcode)
	}

	if len(quoteRequest.Dispatchers) != 1 {
		t.Fatalf("Expected 1 dispatcher, got %d", len(quoteRequest.Dispatchers))
	}

	if quoteRequest.Dispatchers[0].TotalPrice != 100 {
		t.Errorf("Expected total price 100, got %.2f", quoteRequest.Dispatchers[0].TotalPrice)
	}

	if len(quoteRequest.Dispatchers[0].Volumes) != 1 {
		t.Fatalf("Expected 1 volume, got %d", len(quoteRequest.Dispatchers[0].Volumes))
	}

	volume := quoteRequest.Dispatchers[0].Volumes[0]
	if volume.UnitaryPrice != 50 {
		t.Errorf("Expected unitary price 50, got %.2f", volume.UnitaryPrice)
	}
}

func TestFreightService_CreateFreightQuoteRequest_InvalidZipcode(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "invalid",
			},
		},
		Volumes: []models.Volume{{Category: 7, Amount: 1, Price: 100}},
	}

	service := NewFreightService(mockRepo, nil).(*FreightService)

	_, err := service.createFreightQuoteRequest(request)

	if err == nil {
		t.Error("Expected error for invalid zipcode, got nil")
	}
}

func TestFreightService_CreateFreightQuoteRequest_InvalidAmount(t *testing.T) {
	setupEnvVars()
	defer cleanupEnvVars()

	mockRepo := &mocks.MockFreightRepository{}

	request := &models.FreightRequest{
		Recipient: models.Recipient{
			Address: models.Address{
				Zipcode: "01311000",
			},
		},
		Volumes: []models.Volume{
			{Category: 7, Amount: 0, Price: 100}, // Amount inválido
		},
	}

	service := NewFreightService(mockRepo, nil).(*FreightService)

	_, err := service.createFreightQuoteRequest(request)

	if err == nil {
		t.Error("Expected error for invalid amount, got nil")
	}

	if !containsString(err.Error(), "Amount value must be bigger then zero") {
		t.Errorf("Expected specific error message, got '%s'", err.Error())
	}
}

func containsString(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 ||
		(len(str) > 0 && len(substr) > 0 && containsSubstring(str, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
