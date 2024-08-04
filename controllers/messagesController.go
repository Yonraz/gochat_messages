package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yonraz/gochat_messages/services"
)

type MessagesController struct {
	msgSrv *services.MessagesService
}

func NewMessagesController(srv *services.MessagesService) *MessagesController {
	return &MessagesController{
		msgSrv: srv,
	}
}

func (c *MessagesController) GetMessages(ctx *gin.Context) {
	var body struct{
		Sender string `json:"sender"`
		Receiver string `json:"receiver"`
	}
	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing sender and receiver properties",
		})
	}

	conversation, err := c.msgSrv.GetConversation(body.Sender, body.Receiver)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not perform operation",
			"details": err,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"conv" : conversation,
	})
}