package services

import (
	"context"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/mcuadros/go-defaults"
	"log"
	"time"
)

type Service struct {
	Redis               db.Redis
	WindowLength        int `default:"60"`
	AllowedRequestCount int `default:"5"`
	Expiration          int `default:"60"`
	RequestTimeout      int `default:"300"`
}

// Option function template for giving options to middleware service
type Option func(*Service)

func NewService(redis db.Redis, opts ...Option) *Service {
	s := &Service{
		Redis: redis,
	}

	//setting defaults
	defaults.SetDefaults(s)

	//Applying all options given for the service
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (service *Service) CheckAndStoreRate(ctx context.Context, reqData *models.ReqData) bool {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(service.RequestTimeout)*time.Second)
	defer cancel()

	now := time.Now().Unix()

	// storing the request in redis
	_, err := service.Redis.StoreToSortedList(ctxWithTimeout, reqData.UserIp, &models.RedisSaveData{
		IPAddr:     reqData.UserIp,
		URL:        reqData.RequestAddress,
		TimeStamp:  now,
		Expiration: service.Expiration,
	})
	if err != nil {
		log.Printf("failed to store rate for user %s: %v", reqData.UserIp, err)
		return false
	}

	//Get the requester data from redis
	if requesterInRedis, errGettingData := service.Redis.GetSortedList(ctxWithTimeout, reqData.UserIp); errGettingData != nil {
		log.Println("error : Redis failed to get the requesters list : ", errGettingData)
		return false
	} else if requesterInRedis == nil {
		return true
	} else {
		requestCounter := 0
		for _, request := range requesterInRedis {
			if request.TimeStamp < now && request.TimeStamp >= now-int64(service.WindowLength)*60 {
				requestCounter++
			}
		}

		if requestCounter >= service.AllowedRequestCount {
			log.Printf("rate limit exceeded for user %s . The allowed request count is %v", reqData.UserIp, service.AllowedRequestCount)
			return false
		}
		return true
	}
}
