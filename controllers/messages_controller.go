package controllers

import (
	"efficient-api/domain/messages"
	"efficient-api/services"
	"efficient-api/utils/error_utils"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

//type messagesInterface interface {
//	Get(*gin.Context)
//}
//type messenger struct {
//	service services.MessageServiceInterface
//}

func getMessageId(msgIdParam string) (int64, error_utils.MessageErr) {
	msgId, msgErr := strconv.ParseInt(msgIdParam, 10, 64)
	if msgErr != nil {
		return 0, error_utils.NewBadRequestError("message id whould be a number")
	}
	return msgId, nil
}

func  GetMessage(c *gin.Context) {
	msgId, err := getMessageId(c.Param("message_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	message, getErr := services.MessagesService.GetMessage(msgId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, message)
}

func CreateMessage(c *gin.Context) {
	var message messages.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		theErr := error_utils.NewBadRequestError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}
	msg, err := services.MessagesService.CreateMessage(message)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, msg)
}