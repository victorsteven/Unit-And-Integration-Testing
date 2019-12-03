package app

import (
	"efficient-api/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	router = gin.Default()
)

func init() {
	//loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func StartApp() {

	dbdriver := os.Getenv("DBDRIVER")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	database := os.Getenv("DATABASE")
	port := os.Getenv("PORT")

	domain.MessageRepo.Initialize(dbdriver, username, password, port, host, database)
	fmt.Println("DATABASE STARTED")

	routes()

	router.Run(":8080")
}
