package services

import (
	"efficient-api/domain/messages"
	"efficient-api/utils/error_utils"
	"time"
)

var (
	MessagesService messageServiceInterface = &messagesService{}
)

type messagesService struct {}

type messageServiceInterface interface {
	GetMessage(int64) (*messages.Message, error_utils.MessageErr)
	CreateMessage(messages.Message) (*messages.Message, error_utils.MessageErr)
}

func (m *messagesService) GetMessage(msgId int64) (*messages.Message, error_utils.MessageErr) {
	message := &messages.Message{Id: msgId}
	if err := message.Get(); err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) CreateMessage(message messages.Message) (*messages.Message, error_utils.MessageErr) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	message.CreatedAt =  time.Now()
	if err := message.Create(); err != nil {
		return nil, err
	}
	return &message, nil
}