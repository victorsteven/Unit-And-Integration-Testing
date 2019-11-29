package domain

import (
	"efficient-api/utils/error_utils"
	"fmt"
)

func NewMockStorage() MessageRepoInterface {
	return &mockStorage{messages: make(map[int64]*Message)}
}

type mockStorage struct {
	messages map[int64]*Message
	lastID int64
}

func (mr *mockStorage) Get(messageId int64) (*Message, error_utils.MessageErr) {
	result := mr.messages[messageId]
	if result == nil {
		return nil, error_utils.NewBadRequestError(fmt.Sprintf("message %d not found", messageId))
	}
	return result, nil
}

func (mr *mockStorage) Create(msg *Message) (*Message, error_utils.MessageErr) {
	current := mr.messages[msg.Id]
	if current != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("message %d already registered", msg.Id))
	}
	mr.messages[mr.lastID] = msg
	return msg, nil
}

func (mr *mockStorage) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
}