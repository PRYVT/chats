package query

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Id         uuid.UUID   `json:"id" binding:"required"`
	Name       string      `json:"name" binding:"required"`
	UserIds    []uuid.UUID `json:"users" binding:"required"`
	ChangeDate time.Time   `json:"creation_date" binding:"required"`
}
