package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetLimit(c *gin.Context) int {
	limitStr := c.Query("limit")
	limit := 100
	if len(strings.TrimSpace(limitStr)) > 0 {
		limit, err := strconv.Atoi(limitStr)
		if limit > 100 || err != nil || limit <= 0 {
			limit = 100
		}
	}
	return limit
}

func GetOffset(c *gin.Context) int {
	offsetStr := c.Query("offset")
	offset := 0
	if len(strings.TrimSpace(offsetStr)) > 0 {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}
	}
	return offset
}

func GetChatIdParam(c *gin.Context) (uuid.UUID, error) {
	userId := c.Param("chatId")

	if len(strings.TrimSpace(userId)) == 0 {
		return uuid.Nil, fmt.Errorf("path param cant be empty or null")
	}
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}
	return userUuid, nil
}
