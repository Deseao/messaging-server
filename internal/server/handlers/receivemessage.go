package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Deseao/messaging-server/internal/server/middleware"
)

type incomingMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func ReceiveMessage(c *gin.Context) {
	//TODO Implement this
	var request incomingMessage
	c.Bind(&request)
	logger := c.MustGet(middleware.LoggerKey).(*zap.Logger)
	logger.Info("Your text arrived", zap.Any("request", request))
	c.JSON(200, nil)
}
