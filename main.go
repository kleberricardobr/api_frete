package main

import (
	"api_frete/config"
	"api_frete/database"
	"api_frete/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	conf := config.LoadConfig()

	db := database.Database{
		Config: conf,
	}
	err := db.OpenPostgres()
	if err != nil {
		log.Fatalf("Failed on load database: %v", err)
	}

	if err := db.RunMigrations(); err != nil {
		log.Fatal("Failed on run migrations:", err)
	}
	defer db.ClosePostgres()

	server := mux.NewRouter()
	handlers.Add(server)

	log.Printf("Listening on :%d", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), server)
}
