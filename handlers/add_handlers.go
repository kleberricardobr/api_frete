package handlers

import (
	"api_frete/interfaces"
	"api_frete/models"
	"api_frete/repositories"
	"api_frete/services"
	"api_frete/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HandlerDependencies struct {
	FreightService interfaces.IFreightService
	MetricsService interfaces.IFreightMetricsInterface
}

func NewHandlerDependencies() *HandlerDependencies {
	freightRepo := repositories.NewFreightRepository()

	return &HandlerDependencies{
		FreightService: services.NewFreightService(freightRepo, &http.Client{}),
		MetricsService: services.NewMetricsService(freightRepo),
	}
}

func (h *HandlerDependencies) calculateFreight(w http.ResponseWriter, r *http.Request) {
	var freightRequest models.FreightRequest

	err := json.NewDecoder(r.Body).Decode(&freightRequest)
	if err != nil {
		utils.AddError(w, "Invalid JSON Format", http.StatusBadRequest)
		return
	}

	err, resp := h.FreightService.SendFreightQuote(&freightRequest)
	if err != nil {
		utils.AddError(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *HandlerDependencies) getMetrics(w http.ResponseWriter, r *http.Request) {
	lastQuotesStr := r.URL.Query().Get("last_quotes")
	lastQuotes := 0

	if lastQuotesStr != "" {
		var err error
		lastQuotes, err = strconv.Atoi(lastQuotesStr)
		if err != nil || lastQuotes < 0 {
			utils.AddError(w, "Invalid last_quotes parameter. Must be a positive integer",
				http.StatusBadRequest)
			return
		}
	}

	stats, err := h.MetricsService.GetMetrics(lastQuotes)
	if err != nil {
		utils.AddError(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func Add(mux *mux.Router) {
	deps := NewHandlerDependencies()

	mux.HandleFunc("/quote", deps.calculateFreight).Methods("POST")
	mux.HandleFunc("/metrics", deps.getMetrics).Methods("GET")
}
