package messages_db

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"log"
	"os"
	"database/sql"
)

var (
	Client *sql.DB //let assign the users db once here
)


//type dbConnection interface {
//
//}
//
//var (
//	DBConnect dbConnection = &DbClient{}
//)
//
//type DbClient struct {
//	DB *sql.DB
//}

//func (d *DbClient) Connect(ctx context.Context) (*sql.Conn, error) {
//	c, err := d.DB.Conn(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return c, nil
//}

//We can use this function both in our test and our actual database connection
//func NewDBConnection(db *sql.DB) dbConnection {
//	return &DbClient{DB :db}
//}

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
	//var connect DbClient
	Client, err = sql.Open("postgres", db_url)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}

func Mine() {
	fmt.Println("this is called in the main")
}
