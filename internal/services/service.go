package services

import (
	"context"
	"log"
	"main/db"
	"main/internal/models"
	"os"
	"strconv"
	"time"
)

type Service struct {
	Redis db.Redis
}

func NewService(redis db.Redis) Service {
	return Service{
		Redis: redis,
	}
}

func (service *Service) CheckAndStoreRate(ctx context.Context, reqData *models.ReqData) bool {
	requestTimeout, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(requestTimeout)*time.Second)
	defer cancel()

	if requesterInRedis, errGettingData := service.Redis.Get(ctxWithTimeout, reqData.UserIp); errGettingData != nil {
		log.Println(errGettingData)
		return false
	}

}
