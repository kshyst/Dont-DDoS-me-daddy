package services

import (
	"context"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"log"
	"os"
	"strconv"
	"time"
)

const windowLength = 1
const allowedRequestCount = 3
const expiration = 300

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

	now := time.Now().Unix()

	// storing the request in redis
	service.Redis.StoreToSortedList(ctxWithTimeout, reqData.UserIp, &models.RedisSaveData{
		IPAddr:     reqData.UserIp,
		URL:        reqData.RequestAddress,
		TimeStamp:  now,
		Expiration: expiration,
	})

	//Get the requester data from redis
	if requesterInRedis, errGettingData := service.Redis.GetSortedList(ctxWithTimeout, reqData.UserIp); errGettingData != nil {
		log.Println("error : Redis failed to get the requesters list : ", errGettingData)
		return false
	} else if requesterInRedis == nil {
		return true
	} else {
		requestCounter := 0
		for _, request := range requesterInRedis {
			if request.TimeStamp < now && request.TimeStamp >= now-windowLength*60 {
				requestCounter++
			}
		}

		if requestCounter >= allowedRequestCount {
			return false
		}
		return true
	}
}
