package integration__tests

import (
	"database/sql"
	"efficient-api/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"time"
)

const (
	queryTruncateMessage = "TRUNCATE TABLE messages;"
	queryInsertMessage  = "INSERT INTO messages(title, body, created_at) VALUES(?, ?, ?);"
	queryGetAllMessages = "SELECT id, title, body, created_at FROM messages;"
)
var (
	dbConn  *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	os.Exit(m.Run())
}

func database() {
	dbDriver := os.Getenv("DBDRIVER_TEST")
	username := os.Getenv("USERNAME_TEST")
	password := os.Getenv("PASSWORD_TEST")
	host := os.Getenv("HOST_TEST")
	database := os.Getenv("DATABASE_TEST")
	port := os.Getenv("PORT_TEST")

	dbConn = domain.MessageRepo.Initialize(dbDriver, username, password, port, host, database)
}

func refreshMessagesTable() error {

	stmt, err := dbConn.Prepare(queryTruncateMessage)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("Error truncating messages table: %s", err)
	}
	return nil
}

func seedOneMessage() (domain.Message, error) {
	msg := domain.Message{
		Title:     "the title",
		Body:      "the body",
		CreatedAt: time.Now(),
	}
	stmt, err := dbConn.Prepare(queryInsertMessage)
	if err != nil {
		panic(err.Error())
	}
	insertResult, createErr := stmt.Exec(msg.Title, msg.Body, msg.CreatedAt)
	if createErr != nil {
		log.Fatalf("Error creating message: %s", createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Error creating message: %s", createErr)
	}
	msg.Id = msgId
	return msg, nil
}

func seedMessages() ([]domain.Message, error) {
	msgs := []domain.Message{
		{
			Title:     "first title",
			Body:      "first body",
			CreatedAt: time.Now(),
		},
		{
			Title:     "second title",
			Body:      "second body",
			CreatedAt: time.Now(),
		},
	}
	stmt, err := dbConn.Prepare(queryInsertMessage)
	if err != nil {
		panic(err.Error())
	}
	for i, _ := range msgs {
		_, createErr := stmt.Exec(msgs[i].Title, msgs[i].Body, msgs[i].CreatedAt)
		if createErr != nil {
			return nil, createErr
		}
	}
	get_stmt, err := dbConn.Prepare(queryGetAllMessages)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := get_stmt.Query()
	if err != nil {
		return nil,  err
	}
	defer rows.Close()

	results := make([]domain.Message, 0)

	for rows.Next() {
		var msg domain.Message
		if getError := rows.Scan(&msg.Id, &msg.Title, &msg.Body, &msg.CreatedAt); getError != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	return results, nil
}

