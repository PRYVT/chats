package controller

import (
	"net/http"

	"github.com/PRYVT/chats/pkg/models/query"
	"github.com/PRYVT/chats/pkg/query/store/repository"
	"github.com/PRYVT/chats/pkg/query/utils"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/PRYVT/utils/pkg/interfaces"
	"github.com/gin-gonic/gin"
)

type ChatController struct {
	ChatRepo     *repository.ChatRepository
	eventHandler interfaces.EventHandler
}

func NewChatController(userRepo *repository.ChatRepository, eventHandler interfaces.EventHandler) *ChatController {
	return &ChatController{ChatRepo: userRepo, eventHandler: eventHandler}
}

func (ctrl *ChatController) GetChat(c *gin.Context) {

	ChatUuid, err := utils.GetChatIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Chat, err := ctrl.ChatRepo.GetChatById(ChatUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Chat)
}

func (ctrl *ChatController) GetChats(c *gin.Context) {

	limit := utils.GetLimit(c)
	offset := utils.GetOffset(c)

	token := auth.GetTokenFromHeader(c)
	userUuid, err := auth.GetUserUuidFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	Chats, err := ctrl.ChatRepo.GetAllChats(limit, offset, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if Chats == nil {
		Chats = []query.ChatReduced{}
	}
	c.JSON(http.StatusOK, Chats)

}
