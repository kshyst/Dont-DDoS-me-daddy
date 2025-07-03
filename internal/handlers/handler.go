package handlers

import (
	"encoding/json"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"io"
	"net/http"
)

type Handler struct {
	service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	reqData := &models.ReqData{}

	// Getting the body data
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	//check the request validation
	rateAcceptation := handler.service.CheckAndStoreRate(r.Context(), reqData)

	if rateAcceptation {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
}
