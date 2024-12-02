package httphandler

import (
	"context"
	"net/http"

	"github.com/PRYVT/chats/pkg/query/httphandler/controller"
	"github.com/PRYVT/utils/pkg/auth"
	ws "github.com/PRYVT/utils/pkg/websocket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer     *http.Server
	router         *gin.Engine
	ChatController *controller.ChatController
	authMiddleware *auth.AuthMiddleware
	wsController   *ws.WSController
}

func NewHttpHandler(c *controller.ChatController, am *auth.AuthMiddleware, wsController *ws.WSController) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5523",
		Handler: r,
	}
	handler := &HttpHandler{
		router:         r,
		httpServer:     srv,
		ChatController: c,
		authMiddleware: am,
		wsController:   wsController,
	}
	handler.RegisterRoutes()
	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.Use(auth.CORSMiddleware())
	h.router.GET("/ws", h.wsController.OnRequest)
	h.router.Use(h.authMiddleware.AuthenticateMiddleware)
	{
		h.router.GET("Chats/:ChatId", h.ChatController.GetChat)
		h.router.GET("Chats/", h.ChatController.GetChats)
	}
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
