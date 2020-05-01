package handlers

import (
	"github.com/Deseao/messaging-server/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type incomingMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func ReceiveMessage(c *gin.Context) {
	logger := c.MustGet(middleware.LoggerKey).(*zap.Logger)
	var request incomingMessage
	err := c.Bind(&request)
	if err != nil {
		logger.Error("Failed to parse request", zap.Error(err))
	} else {
		logger.Info("Your text arrived", zap.Any("request", request))
	}
	c.JSON(200, nil)
}
