package Daddy

import (
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
)

func WithWindowLength(seconds int) services.Option {
	return func(s *services.Service) {
		s.WindowLength = seconds
	}
}

func WithAllowedRequestCount(count int) services.Option {
	return func(s *services.Service) {
		s.AllowedRequestCount = count
	}
}

func WithExpiration(seconds int) services.Option {
	return func(s *services.Service) {
		s.Expiration = seconds
	}
}

func WithRequestTimeout(seconds int) services.Option {
	return func(s *services.Service) {
		s.RequestTimeout = seconds
	}
}
