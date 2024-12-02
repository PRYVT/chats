package main

import (
	"os"

	"github.com/PRYVT/chats/pkg/command/httphandler"
	"github.com/PRYVT/chats/pkg/command/httphandler/controller"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	aut := auth.NewAuthMiddleware()
	uc := controller.NewChatsController()
	h := httphandler.NewHttpHandler(uc, aut)

	h.Start()
}
