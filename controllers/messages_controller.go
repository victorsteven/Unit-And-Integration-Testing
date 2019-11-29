package controllers

import (
	"efficient-api/domain"
	"efficient-api/services"
	"efficient-api/utils/error_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		return 0, error_utils.NewBadRequestError("message id should be a number")
	}
	return msgId, nil
}

func GetMessage(c *gin.Context) {
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
	var message domain.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		theErr := error_utils.NewUnprocessibleEntityError("invalid json body")
		c.JSON(theErr.Status(), theErr)
		return
	}
	msg, err := services.MessagesService.CreateMessage(&message)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, msg)
}

//func Ping(c *gin.Context) {
//	c.JSON(http.StatusOK, "This is the route")
//}
