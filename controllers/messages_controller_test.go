package controllers

import (
	"bytes"
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
	updateMessageService func(message *domain.Message) (*domain.Message, error_utils.MessageErr)
	deleteMessageService func(msgId int64) error_utils.MessageErr
	getAllMessageService func() ([]domain.Message, error_utils.MessageErr)
)

type serviceMock struct {}

func (sm *serviceMock) GetMessage(msgId int64) (*domain.Message, error_utils.MessageErr) {
	return getMessageService(msgId)
}
func (sm *serviceMock) CreateMessage(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
	return createMessageService(message)
}
func (sm *serviceMock) UpdateMessage(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
	return updateMessageService(message)
}
func (sm *serviceMock) DeleteMessage(msgId int64) error_utils.MessageErr {
	return deleteMessageService(msgId)
}
func (sm *serviceMock) GetAllMessages() ([]domain.Message, error_utils.MessageErr) {
	return getAllMessageService()
}

///////////////////////////////////////////////////////////////
// Start of "GetMessage" test cases
///////////////////////////////////////////////////////////////
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
///////////////////////////////////////////////////////////////
// End of "GetMessage" test cases
///////////////////////////////////////////////////////////////



///////////////////////////////////////////////////////////////
// Start of "CreateMessage" test cases
///////////////////////////////////////////////////////////////
func TestCreateMessage_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	createMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:        1,
			Title:     "the title",
			Body:      "the body",
		}, nil
	}
	jsonBody := `{"title": "the title", "body": "the body"}`
	r := gin.Default()
	req, err := http.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.POST("/messages", CreateMessage)
	r.ServeHTTP(rr, req)

	var message domain.Message
	err = json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "the title", message.Title)
	assert.EqualValues(t, "the body", message.Body)
}

func TestCreateMessage_Invalid_Json(t *testing.T) {
	inputJson := `{"title": 1234, "body": "the body"}`
	r := gin.Default()
	req, err := http.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.POST("/messages", CreateMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

//This test is not really necessary here, because it has been handled in the service test
func TestCreateMessage_Empty_Body(t *testing.T) {
	services.MessagesService = &serviceMock{}
	createMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewUnprocessibleEntityError("Please enter a valid body")
	}
	inputJson := `{"title": "the title", "body": ""}`
	r := gin.Default()
	req, err := http.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.POST("/messages", CreateMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}
//This test is not really necessary here, because it has been handled in the service test
func TestCreateMessage_Empty_Title(t *testing.T) {
	services.MessagesService = &serviceMock{}
	createMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewUnprocessibleEntityError("Please enter a valid title")
	}
	inputJson := `{"title": "", "body": "the body"}`
	r := gin.Default()
	req, err := http.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.POST("/messages", CreateMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid title", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}
///////////////////////////////////////////////////////////////
// End of "CreateMessage" test cases
///////////////////////////////////////////////////////////////



///////////////////////////////////////////////////////////////
// Start of "UpdateMessage" test cases
///////////////////////////////////////////////////////////////
func TestUpdateMessage_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	updateMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:        1,
			Title:     "update title",
			Body:      "update body",
		}, nil
	}
	jsonBody := `{"title": "update title", "body": "update body"}`
	r := gin.Default()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/messages/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.PUT("/messages/:message_id", UpdateMessage)
	r.ServeHTTP(rr, req)

	var message domain.Message
	err = json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "update title", message.Title)
	assert.EqualValues(t, "update body", message.Body)
}

//We dont need to mock the service method here, because we wont call it
func TestUpdateMessage_Invalid_Id(t *testing.T) {
	jsonBody := `{"title": "update title", "body": "update body"}`
	r := gin.Default()
	id := "abc"
	req, err := http.NewRequest(http.MethodPut, "/messages/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.PUT("/messages/:message_id", UpdateMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "message id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

//When for instance an integer is provided instead of a string
func TestUpdateMessage_Invalid_Json(t *testing.T) {
	inputJson := `{"title": 1234, "body": "the body"}`
	r := gin.Default()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/messages/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.PUT("/messages/:message_id", UpdateMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}


//This test is not really necessary here, because it has been handled in the service test
func TestUpdateMessage_Empty_Body(t *testing.T) {
	services.MessagesService = &serviceMock{}
	updateMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewUnprocessibleEntityError("Please enter a valid body")
	}
	inputJson := `{"title": "the title", "body": ""}`
	id := "1"
	r := gin.Default()
	req, err := http.NewRequest(http.MethodPut, "/messages/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.PUT("/messages/:message_id", UpdateMessage)
	r.ServeHTTP(rr, req)
	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}
//This test is not really necessary here, because it has been handled in the service test
func TestUpdateMessage_Empty_Title(t *testing.T) {
	services.MessagesService = &serviceMock{}
	updateMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewUnprocessibleEntityError("Please enter a valid title")
	}
	inputJson := `{"title": "", "body": "the body"}`
	id := "1"
	r := gin.Default()
	req, err := http.NewRequest(http.MethodPut, "/messages/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.PUT("/messages/:message_id", UpdateMessage)
	r.ServeHTTP(rr, req)
	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid title", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

//Other errors can happen when we try to update the message
func TestUpdateMessage_Error_Updating(t *testing.T) {
	services.MessagesService = &serviceMock{}
	updateMessageService = func(message *domain.Message) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewInternalServerError("error when updating message")
	}
	jsonBody := `{"title": "update title", "body": "update body"}`
	r := gin.Default()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/messages/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.PUT("/messages/:message_id", UpdateMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)

	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error when updating message", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}
///////////////////////////////////////////////////////////////
// End of "UpdateMessage" test cases
///////////////////////////////////////////////////////////////



///////////////////////////////////////////////////////////////
// Start of "DeleteMessage" test cases
///////////////////////////////////////////////////////////////
func TestDeleteMessage_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	deleteMessageService = func(msg int64) error_utils.MessageErr {
		return nil
	}
	r := gin.Default()
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/messages/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.DELETE("/messages/:message_id", DeleteMessage)
	r.ServeHTTP(rr, req)

	var response = make(map[string]string)
	theErr := json.Unmarshal(rr.Body.Bytes(), &response)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, response["status"], "deleted")
}

//We wont call the service Delete method here, so no need to mock it
func TestDeleteMessage_Invalid_Id(t *testing.T) {

	r := gin.Default()
	id := "abc"
	req, err := http.NewRequest(http.MethodDelete, "/messages/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.DELETE("/messages/:message_id", DeleteMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "message id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

//If for any reason, our update didnt happen(e.g server error, etc), This is an error from the service, but the controller conditions where met.
//Maybe the message does not exist, or the server timeout
func TestDeleteMessage_Failure(t *testing.T) {
	services.MessagesService = &serviceMock{}
	deleteMessageService = func(msg int64) error_utils.MessageErr {
		return error_utils.NewInternalServerError("error deleting message")
	}
	r := gin.Default()
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/messages/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.DELETE("/messages/:message_id", DeleteMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error deleting message", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}
///////////////////////////////////////////////////////////////
// End of "DeleteMessage" test cases
///////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////
// Start of "GetAllMessages" test cases
///////////////////////////////////////////////////////////////
func TestGetAllMessages_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getAllMessageService = func() ([]domain.Message, error_utils.MessageErr) {
		 return []domain.Message{
			{
				Id:        1,
				Title:     "first title",
				Body:      "first body",
			},
			{
				Id:        2,
				Title:     "second title",
				Body:      "second body",
			},
		}, nil
	}
	r := gin.Default()
	req, err := http.NewRequest(http.MethodGet, "/messages", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.GET("/messages", GetAllMessages)
	r.ServeHTTP(rr, req)

	var messages []domain.Message
	theErr := json.Unmarshal(rr.Body.Bytes(), &messages)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, messages)
	assert.EqualValues(t, messages[0].Id, 1)
	assert.EqualValues(t, messages[0].Title, "first title")
	assert.EqualValues(t, messages[0].Body, "first body")
	assert.EqualValues(t, messages[1].Id, 2)
	assert.EqualValues(t, messages[1].Title, "second title")
	assert.EqualValues(t, messages[1].Body, "second body")
}

//For any reason we could not get the messages
func TestGetAllMessages_Failure(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getAllMessageService = func() ([]domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewInternalServerError("error getting messages")
	}
	r := gin.Default()
	req, err := http.NewRequest(http.MethodGet, "/messages", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.GET("/messages", GetAllMessages)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error getting messages", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}
