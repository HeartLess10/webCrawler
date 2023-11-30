package main

import (
	"log"
	"os"
	"webScraping/databases"
	"webScraping/server"

	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "Crawler-App-Prefix ", log.LstdFlags)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error: Failed to load .env, ", err)
	}
	//testing mongodb
	databases.MongoDbTestComblixQuery(os.Getenv("DB_CLUSTER_PATH"), os.Getenv("DATA_BASE_NAME"), os.Getenv("DATA_BASE_COLLECTION_NAME"), logger)

	server.StartServer(logger)

}
