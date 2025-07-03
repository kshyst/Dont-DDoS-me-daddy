package main

import (
	"github.com/joho/godotenv"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/handlers"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"log"
	"net/http"
	"os"
)

func main() {

	userRedis, err := db.NewRedis()
	if err != nil {
		log.Fatal("redis failed to start")
	}
	userSrv := services.NewService(userRedis)
	userHandler := handlers.NewHandler(userSrv)

	if envLoadingError := godotenv.Load(); envLoadingError != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/req", userHandler.Handle)

	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	err = http.ListenAndServe(address+":"+port, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}
