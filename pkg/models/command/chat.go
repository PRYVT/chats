package command

import "github.com/google/uuid"

type CreateChat struct {
	Id      string             `json:"id" binding:"required"`
	Name    string             `json:"name" binding:"required"`
	UserIds map[uuid.UUID]bool `json:"users" binding:"required"`
}
