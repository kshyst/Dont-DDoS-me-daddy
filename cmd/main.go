package main

import (
	Daddy "github.com/kshyst/Dont-DDoS-me-daddy"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/handlers"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"log"
	"net/http"
)

func main() {
	// Load yaml configs
	config, loadConfigErr := models.LoadConfig("configs.yaml")
	if loadConfigErr != nil {
		log.Fatalf("Error loading configs.yaml : %s", loadConfigErr)
	}

	// Instantiating redis
	userRedis, err := db.NewRedis()
	if err != nil {
		log.Fatal("redis failed to start")
	}

	// Create service and the handler
	userSrv := services.NewService(
		userRedis,
		Daddy.WithExpiration(models.GetIntOfStrings(config.Options.RedisExpiration)),
		Daddy.WithAllowedRequestCount(models.GetIntOfStrings(config.Options.AllowedRequestCount)),
		Daddy.WithWindowLength(models.GetIntOfStrings(config.Options.WindowLength)),
		Daddy.WithRequestTimeout(models.GetIntOfStrings(config.Options.ContextTimeout)),
	)

	userHandler := handlers.NewHandler(userSrv)

	// Endpoint for rate limiting
	http.HandleFunc("/req", userHandler.Handle)

	// Start the rate limiter server
	err = http.ListenAndServe(config.Server.Address+":"+config.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Rate Limiter Server Started")
}
