package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const LoggerKey string = "logger"

func SetupLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestLog := logger.With(zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		c.Set(LoggerKey, requestLog)
	}
}

func AddToContext(key string, value interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, value)
	}
}
