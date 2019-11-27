package messages

import (
	"efficient-api/database/postgres/messages_db"
	"efficient-api/utils/error_utils"
	"fmt"
)

const (
	queryGetMessage = "SELECT id, title, body, created_at FROM messages WHERE id=$1;"
	queryInsertMessage = "INSERT INTO messages(title, body, created_at) VALUES($1, $2, $3) RETURNING id;"
)

func (m *Message) Get() error_utils.MessageErr {
	stmt, err := messages_db.Client.Prepare(queryGetMessage)
	if err != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to prepare message: %s", err.Error()))
	}
	defer stmt.Close()

	result := stmt.QueryRow(m.Id)

	if getError := result.Scan(&m.Id, &m.Title, &m.Body, &m.CreatedAt); getError != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
	}
	return nil
}

func (m *Message) Create() error_utils.MessageErr {
	stmt, err := messages_db.Client.Prepare(queryInsertMessage)
	if err != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}
	defer stmt.Close()

	var messageId int64
	createErr := stmt.QueryRow(m.Title, m.Body, m.CreatedAt).Scan(&messageId)
	if createErr != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", createErr.Error()))
	}
	m.Id = messageId
	return nil
}