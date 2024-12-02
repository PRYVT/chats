package controller

import (
	"net/http"

	"github.com/PRYVT/chats/pkg/aggregates"
	"github.com/PRYVT/chats/pkg/models/command"
	"github.com/PRYVT/chats/pkg/query/utils"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatsController struct {
}

func NewChatsController() *ChatsController {
	return &ChatsController{}
}

func (ctrl *ChatsController) CreateChat(c *gin.Context) {

	token := auth.GetTokenFromHeader(c)
	userUuid, err := auth.GetUserUuidFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var m command.CreateChat
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chatUuid := uuid.MustParse(m.Id)
	ua, err := aggregates.NewChatAggregate(chatUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	m.UserIds[userUuid] = true
	err = ua.CreateChat(m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (ctrl *ChatsController) AddChatMessage(c *gin.Context) {

	token := auth.GetTokenFromHeader(c)
	userUuid, err := auth.GetUserUuidFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	chatUuid, err := utils.GetChatIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var m command.AddChatMessage
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ua, err := aggregates.NewChatAggregate(chatUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = ua.AddChatMessage(m, userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
