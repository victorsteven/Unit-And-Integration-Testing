package services

import (
	"efficient-api/domain"
	"efficient-api/utils/error_utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

//type methodTypes struct {
//	getMessage func(messageId int64) (*domain.Message, error_utils.MessageErr)
//	createMessage func(msg *domain.Message) (*domain.Message, error_utils.MessageErr)
//}
//var theUtil methodTypes

var (
	getMessage func(messageId int64) (*domain.Message, error_utils.MessageErr)
	createMessage func(msg *domain.Message) (*domain.Message, error_utils.MessageErr)
)

type getDBMock struct {}

func (m *getDBMock) Get(messageId int64) (*domain.Message, error_utils.MessageErr){
	//return theUtil.getMessage(messageId)
	return getMessage(messageId)

}
func (m *getDBMock) Create(msg *domain.Message) (*domain.Message, error_utils.MessageErr){
	return createMessage(msg)
	//return theUtil.createMessage(msg)

}
func (m *getDBMock) Initialize(string, string, string, string, string, string){}

func TestMessagesService_GetMessage_Success(t *testing.T) {
	tm := time.Now()
	domain.MessageRepo = &getDBMock{}

	getMessage = func(messageId int64) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:        1,
			Title:     "the title",
			Body:      "the body",
			CreatedAt: tm,
		}, nil
	}
	msg, err := MessagesService.GetMessage(1)
	fmt.Println("this is the message: ", msg)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, "the title", msg.Title)
	assert.EqualValues(t, "the body", msg.Body)
	assert.EqualValues(t, tm, msg.CreatedAt)

	//tests := []struct {
	//	getMessage
	//	msgId int64
	//}{
	//	{
	//		msgId: 1,
	//		getMessage:  func(messageId int64) (*domain.Message, error_utils.MessageErr) {
	//			return &domain.Message{
	//				Id:        1,
	//				Title:     "the title",
	//				Body:      "the body",
	//				CreatedAt: tm,
	//			}, nil
	//		},
	//	},
	//}
	//for _, tt := range tests {
	//
	//	fmt.Println("this is the type of the variable: ", reflect.TypeOf(tt.msgId))
	//
	//	msg, err := MessagesService.GetMessage(tt.msgId)
	//
	//	domain.MessageRepo = &getDBMock{}
	//
	//	assert.NotNil(t, msg)
	//	assert.Nil(t, err)
	//	assert.EqualValues(t, 1, msg.Id)
	//	assert.EqualValues(t, "the title", msg.Title)
	//	assert.EqualValues(t, "the body", msg.Body)
	//	assert.EqualValues(t, tm, msg.CreatedAt)
	//}
}

func TestMessagesService_GetMessage_NotFoundID(t *testing.T) {
	//tm := time.Now()
	domain.MessageRepo = &getDBMock{}

	getMessage = func(messageId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("the id is not found")
	}
	msg, err := MessagesService.GetMessage(123)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "the id is not found", err.Message())
	assert.EqualValues(t, "not_found", err.Error())
}