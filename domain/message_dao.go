package domain

import (
	"efficient-api/utils/error_utils"
	_ "github.com/lib/pq"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"database/sql"
)

var (
	MessageRepo messageRepoInterface = &messageRepo{}
)

const (
	queryGetMessage = "SELECT id, title, body, created_at FROM messages WHERE id=$1;"
	queryInsertMessage = "INSERT INTO messages(title, body, created_at) VALUES($1, $2, $3) RETURNING id;"
)
type messageRepoInterface interface {
	Get(int64) (*Message, error_utils.MessageErr)
	Create(*Message) (*Message, error_utils.MessageErr)
	Initialize(string, string, string, string, string, string)
}
type messageRepo struct {
	db *sql.DB
}

func (mr *messageRepo) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	mr.db, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error connecting to postgres:", err)
	} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
	}
}

func NewMessageRepository(db *sql.DB) messageRepoInterface {
	return &messageRepo{db: db}
}

func (mr *messageRepo) Get(messageId int64) (*Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryGetMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to prepare message: %s", err.Error()))
	}
	defer stmt.Close()

	var msg Message
	result := stmt.QueryRow(messageId)
	if getError := result.Scan(&msg.Id, &msg.Title, &msg.Body, &msg.CreatedAt); getError != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
	}
	return &msg, nil
}

func (mr *messageRepo) Create(msg *Message) (*Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryInsertMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}
	defer stmt.Close()

	var messageId int64
	createErr := stmt.QueryRow(msg.Title, msg.Body, msg.CreatedAt).Scan(&messageId)
	if createErr != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", createErr.Error()))
	}
	msg.Id = messageId
	return msg, nil
}