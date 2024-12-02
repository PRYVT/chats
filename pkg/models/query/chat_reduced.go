package query

import (
	"github.com/google/uuid"
)

type ChatReduced struct {
	Id   uuid.UUID `json:"id" binding:"required"`
	Name string    `json:"name" binding:"required"`
}