package main

import (
	"fmt"
	"os"

	"github.com/Valgard/godotenv"
	"github.com/gofiber/fiber/v2/log"

	"github.com/sentinel/geospatial/storage"
)

func main() {

	err := godotenv.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbConnStr := os.Getenv("DB_CONN_STR")
	if dbConnStr == "" {
		err1 := fmt.Errorf("failed to load DB_CONN_STR from .env\n	check if variable exists")
		log.Fatal(err1)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		err1 := fmt.Errorf("failed to load DB_NAME from .env\n	check if variable exists")
		log.Fatal(err1)
	}

	mongoDBClient, err := storage.MongoDBInit(dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	ambulanceCollection := mongoDBClient.Database(dbName).Collection("ambulances")

	ambulanceStorage, err := storage.NewMongoDBAmbulanceStorage(ambulanceCollection)
	if err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":3000", ambulanceStorage)

	//	start server and listen
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
