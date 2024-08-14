package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yonraz/gochat_messages/services"
)

type MessagesController struct {
	msgSrv services.MessagesServiceInterface
}
type GetMessageReqBody struct {
	Sender string 
	Receiver string 
}

func NewMessagesController(srv services.MessagesServiceInterface) *MessagesController {
	return &MessagesController{
		msgSrv: srv,
	}
}

func (c *MessagesController) GetMessages(ctx *gin.Context) {
	sender, senderExists := ctx.GetQuery("sender")
	receiver, recExists := ctx.GetQuery("receiver")
	offsetQuery := ctx.DefaultQuery("offset", "0")
	offset, queryErr := strconv.Atoi(offsetQuery)
	if queryErr != nil {
		log.Println("page query invalid, defaulting to 0.")
		offsetQuery = "0"
		offset = 0
	}
	

	if !senderExists || !recExists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing sender and receiver query params",
		})
		return
	}
	log.Printf("request to get messages with sender %v and receiver %v\n", sender, receiver)

	conversation, err := c.msgSrv.GetConversationWithMessages(sender, receiver, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not perform operation",
			"details": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages" : conversation.Messages,
	})
}