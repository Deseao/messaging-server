package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Deseao/messaging-server/internal/server/handlers"
	"github.com/Deseao/messaging-server/internal/server/middleware"
	"github.com/Deseao/messaging-server/internal/state"
)

type Config struct {
	Production bool
	Logger     *zap.Logger
	Accept     string
	Port       string
	State      *state.State
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func New(cfg *Config) *Server {
	if cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	routes := gin.New()
	//TODO: Gin.Recovery
	routes.HandleMethodNotAllowed = true
	routes.Use(middleware.SetupLogger(cfg.Logger))
	routes.Use(middleware.AddToContext("state", cfg.State))
	routes.GET("/ping", healthHandler)
	routes.POST("/signup", handlers.Signup)
	routes.POST("/sms", handlers.ReceiveMessage)
	return &Server{
		server: &http.Server{
			Addr:           cfg.Accept + ":" + cfg.Port,
			Handler:        routes,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		logger: cfg.Logger,
	}
}

func (serv *Server) Run(secure bool, certFile string, keyFile string) {
	if secure {
		serv.logger.Info("starting HTTPS server...", zap.String("certFile", certFile), zap.String("keyFile", keyFile))
		if err := serv.server.ListenAndServeTLS(certFile, keyFile); !errors.Is(err, http.ErrServerClosed) {
			serv.logger.Fatal("service shutdown unexpectedly", zap.Error(err))
		}
	} else {
		serv.logger.Info("starting HTTP server...")
		if err := serv.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serv.logger.Fatal("service shutdown unexpectedly", zap.Error(err))
		}
	}
}

func (serv *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := serv.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
