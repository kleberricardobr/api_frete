package config

import (
	"api_frete/models"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig() (config *models.ConfigModel) {
	config = &models.ConfigModel{
		Port: 8080,
		Db: models.DatabaseModel{
			Port: 5432,
			Host: "127.0.0.1",
		},
	}

	godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432
	}

	config.Port = port
	config.Db.Port = dbPort
	config.Db.User = os.Getenv("DB_USER")
	config.Db.Pass = os.Getenv("DB_PASS")
	config.Db.Host = os.Getenv("DB_HOST")
	config.Db.Name = os.Getenv("DB_NAME")

	return
}
