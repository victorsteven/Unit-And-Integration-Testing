package domain

import (
	"database/sql"
	"efficient-api/utils/error_utils"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	MessageRepo MessageRepoInterface = &messageRepo{}
)

const (
	queryGetMessage    = "SELECT id, title, body, created_at FROM messages WHERE id=?;"
	queryInsertMessage = "INSERT INTO messages(title, body, created_at) VALUES(?, ?, ?);"
	queryUpdateMessage = "UPDATE messages SET title=?, body=? WHERE id=?;"
	queryDeleteMessage = "DELETE FROM messages WHERE id=?;"
)

type MessageRepoInterface interface {
	Get(int64) (*Message, error_utils.MessageErr)
	Create(*Message) (*Message, error_utils.MessageErr)
	Update(*Message) (*Message, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	Initialize(string, string, string, string, string, string)
}
type messageRepo struct {
	db *sql.DB
}

func (mr *messageRepo) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	mr.db, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error connecting to postgres:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}
}

func NewMessageRepository(db *sql.DB) MessageRepoInterface {
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

	insertResult, createErr := stmt.Exec(msg.Title, msg.Body, msg.CreatedAt)
	if createErr != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", createErr.Error()))
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

	result, updateErr := stmt.Exec(msg.Title, msg.Body, msg.Id)
	fmt.Println("this is the result: ", result)
	if updateErr != nil {
		return nil, error_utils.NewInternalServerError("error when trying to update message first")
	}
	return msg, nil
}

func (mr *messageRepo) Delete(msgId int64) error_utils.MessageErr {
	stmt, err := mr.db.Prepare(queryDeleteMessage)
	if err != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to delete message: %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(msgId); err != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to delete message %s", err.Error()))
	}
	return nil
}
