package messages_db

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"log"
	"os"
)

import "database/sql"

type dbConnection interface {

}

type dbClient struct {
	db *sql.DB
}

//We can use this function both in our test and our actual database connection
func NewDBConnection(db *sql.DB) dbConnection {
	return &dbClient{db:db}
}

func init(){
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	database := os.Getenv("DATABASE")
	port := os.Getenv("PORT")

	db_url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)

	var err error
	var connect dbClient
	connect.db, err = sql.Open("postgres", db_url)
	if err != nil {
		panic(err)
	}
	if err = connect.db.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}

func Mine() {
	fmt.Println("this is called in the main")
}
