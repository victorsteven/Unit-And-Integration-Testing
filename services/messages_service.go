package services

import (
	"efficient-api/domain"
	"efficient-api/utils/error_utils"
	"time"
)

var (
	MessagesService messageServiceInterface = &messagesService{}
)

type messagesService struct {}

type messageServiceInterface interface {
	GetMessage(int64) (*domain.Message, error_utils.MessageErr)
	CreateMessage(*domain.Message) (*domain.Message, error_utils.MessageErr)
}

func (m *messagesService) GetMessage(msgId int64) (*domain.Message, error_utils.MessageErr) {
	message, err := domain.MessageRepo.Get(msgId);
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) CreateMessage(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	message.CreatedAt =  time.Now()
	message, err := domain.MessageRepo.Create(message);
	if err != nil {
		return nil, err
	}
	return message, nil
}