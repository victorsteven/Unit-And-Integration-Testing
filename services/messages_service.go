package services

import (
	"efficient-api/domain"
	"efficient-api/utils/error_utils"
	"time"
)

var (
	MessagesService messageServiceInterface = &messagesService{}
)

type messagesService struct{}

type messageServiceInterface interface {
	GetMessage(int64) (*domain.Message, error_utils.MessageErr)
	CreateMessage(*domain.Message) (*domain.Message, error_utils.MessageErr)
	UpdateMessage(*domain.Message) (*domain.Message, error_utils.MessageErr)
	DeleteMessage(int64) error_utils.MessageErr
}

func (m *messagesService) GetMessage(msgId int64) (*domain.Message, error_utils.MessageErr) {
	message, err := domain.MessageRepo.Get(msgId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) CreateMessage(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	message.CreatedAt = time.Now()
	message, err := domain.MessageRepo.Create(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) UpdateMessage(message *domain.Message) (*domain.Message, error_utils.MessageErr) {

	current, err := m.GetMessage(message.Id)
	if err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	current.Title = message.Title
	current.Body = message.Body

	updateMsg, err := domain.MessageRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updateMsg, nil
}

func (m *messagesService) DeleteMessage(msgId int64) error_utils.MessageErr {
	_, err := m.GetMessage(msgId)
	if err != nil {
		return err
	}
	return domain.MessageRepo.Delete(msgId)
}