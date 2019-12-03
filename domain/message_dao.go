package domain

import (
	"database/sql"
	"efficient-api/utils/error_formats"
	"efficient-api/utils/error_utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	MessageRepo messageRepoInterface = &messageRepo{}
)

const (
	queryGetMessage    = "SELECT id, title, body, created_at FROM messages WHERE id=?;"
	queryInsertMessage = "INSERT INTO messages(title, body, created_at) VALUES(?, ?, ?);"
	queryUpdateMessage = "UPDATE messages SET title=?, body=? WHERE id=?;"
	queryDeleteMessage = "DELETE FROM messages WHERE id=?;"
	queryGetAllMessages = "SELECT id, title, body, created_at FROM messages;"
)

type messageRepoInterface interface {
	Get(int64) (*Message, error_utils.MessageErr)
	Create(*Message) (*Message, error_utils.MessageErr)
	Update(*Message) (*Message, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	GetAll() ([]Message, error_utils.MessageErr)
	Initialize(string, string, string, string, string, string) *sql.DB
}
type messageRepo struct {
	db *sql.DB
}

func (mr *messageRepo) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *sql.DB  {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	mr.db, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		log.Fatal("This is the error connecting to the database:", err)
	}
	fmt.Printf("We are connected to the %s database", Dbdriver)

	return mr.db
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
		fmt.Println("this is the error man: ", getError)
		return nil,  error_formats.ParseError(getError)
	}
	return &msg, nil
}

func (mr *messageRepo) GetAll() ([]Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryGetAllMessages)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to prepare all messages: %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil,  error_formats.ParseError(err)
	}
	defer rows.Close()

	results := make([]Message, 0)

	for rows.Next() {
		var msg Message
		if getError := rows.Scan(&msg.Id, &msg.Title, &msg.Body, &msg.CreatedAt); getError != nil {
			return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
		}
		results = append(results, msg)
	}
	if len(results) == 0 {
		return nil, error_utils.NewNotFoundError("no records found")
	}
	return results, nil
}

func (mr *messageRepo) Create(msg *Message) (*Message, error_utils.MessageErr) {
	fmt.Println("WE REACHED THE DOMAIN")
	stmt, err := mr.db.Prepare(queryInsertMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}
	fmt.Println("WE DIDNT REACH HERE")

	defer stmt.Close()

	insertResult, createErr := stmt.Exec(msg.Title, msg.Body, msg.CreatedAt)
	if createErr != nil {
		return nil,  error_formats.ParseError(createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}
	msg.Id = msgId

	return msg, nil
}

func (mr *messageRepo) Update(msg *Message) (*Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryUpdateMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to update: %s", err.Error()))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(msg.Title, msg.Body, msg.Id)
	if updateErr != nil {
		return nil,  error_formats.ParseError(updateErr)
	}
	return msg, nil
}

func (mr *messageRepo) Delete(msgId int64) error_utils.MessageErr {
	stmt, err := mr.db.Prepare(queryDeleteMessage)
	if err != nil {
		fmt.Println("this is the delete error: ", err)
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to delete message: %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(msgId); err != nil {
		fmt.Println("this is the second delete error: ", err)
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to delete message %s", err.Error()))
	}
	return nil
}
