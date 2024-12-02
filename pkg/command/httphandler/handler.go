package httphandler

import (
	"context"
	"net/http"

	"github.com/PRYVT/chats/pkg/command/httphandler/controller"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer     *http.Server
	router         *gin.Engine
	ChatController *controller.ChatsController
	authMiddleware *auth.AuthMiddleware
}

func NewHttpHandler(c *controller.ChatsController, m *auth.AuthMiddleware) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5522",
		Handler: r,
	}
	handler := &HttpHandler{
		router:         r,
		httpServer:     srv,
		ChatController: c,
		authMiddleware: m,
	}

	handler.RegisterRoutes()

	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.Use(auth.CORSMiddleware())
	h.router.Use(h.authMiddleware.AuthenticateMiddleware)
	h.router.POST("chats/", h.ChatController.CreateChat)
	h.router.POST("chats/:chatId/message", h.ChatController.AddChatMessage)
}

func (h *HttpHandler) Start() error {
	return h.httpServer.ListenAndServe()
}

func (h *HttpHandler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Error during reading response body")
	}
}
