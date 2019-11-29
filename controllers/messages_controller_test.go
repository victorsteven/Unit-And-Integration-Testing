package controllers

import (
	"efficient-api/domain"
	"efficient-api/services"
	"efficient-api/utils/error_utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	tm = time.Now()
	getMessageService func(msgId int64) (*domain.Message, error_utils.MessageErr)
	createMessageService func(message *domain.Message) (*domain.Message, error_utils.MessageErr)
)

type serviceMock struct {}

func (sm *serviceMock) GetMessage(msgId int64) (*domain.Message, error_utils.MessageErr) {
	return getMessageService(msgId)
}
func (sm *serviceMock) CreateMessage(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
	return createMessageService(message)
}

func TestGetMessage_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getMessageService = func(msgId int64) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:        1,
			Title:     "the title",
			Body:      "the body",
		}, nil
	}
	msgId := "1" //this has to be a string, because is passed through the url
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	var message domain.Message
	err := json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "the title", message.Title)
	assert.EqualValues(t, "the body", message.Body)

}

//When an invalid id id passed. No need to mock the service here because we will never call it
func TestGetMessage_Invalid_Id(t *testing.T) {
	msgId := "abc" //this has to be a string, because is passed through the url
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "message id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

//We will call the service method here, so we need to mock it
func TestGetMessage_Message_Not_Found(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getMessageService = func(msgId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("message not found")
	}
	msgId := "1" //valid id
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "message not found", apiErr.Message())
	assert.EqualValues(t, "not_found", apiErr.Error())

}

//We will call the service method here, so we need to mock it
//If for any reason, we could not get the message
func TestGetMessage_Message_Database_Error(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getMessageService = func(msgId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewInternalServerError("database error")
	}
	msgId := "1" //valid id
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "database error", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

func TestCreateMessage(t *testing.T) {

}