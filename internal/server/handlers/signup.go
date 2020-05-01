package handlers

import (
	"net/http"

	"github.com/Deseao/messaging-server/internal/server/middleware"
	"github.com/Deseao/messaging-server/internal/state"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignupPayload struct {
	Number   string `json:"number"`
	Timezone string `json:"timezone"`
}

func Signup(c *gin.Context) {
	state := c.MustGet("state").(*state.State)
	logger := c.MustGet(middleware.LoggerKey).(*zap.Logger)
	var payload SignupPayload
	err := c.Bind(&payload)
	if err != nil {
		logger.Error("Failed to bind signup payload", zap.Error(err))
	}
	state.AddSubscriber(payload.Number)
	logger.Info("Signup completed", zap.Any("subscribers", state.Subscribers))
	c.JSON(http.StatusNotImplemented, nil)
}
