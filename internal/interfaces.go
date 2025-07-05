package internal

import (
	"context"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"net/http"
)

type Service interface {
	CheckAndStoreRate(ctx context.Context, reqData *models.ReqData) bool
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
