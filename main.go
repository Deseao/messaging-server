package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deseao/messaging-server/internal/config"
	"github.com/Deseao/messaging-server/internal/messaging"
	"github.com/Deseao/messaging-server/internal/server"
	"github.com/SpiderOak/errstack"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.InitConfig("messaging-server")
	if err != nil {
		panic("failed to initialize config: " + err.Error())
	}
	logger, err := makeLogger(cfg.Logger.Production, cfg.Logger.Level)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync() //nolint //Don't care about error from this since it happens even when things are fine
	logger.Info("logger initialized", zap.Bool("production", cfg.Logger.Production), zap.String("level", cfg.Logger.Level))

	logger.Info("initializing server...", zap.String("acceptAddr", cfg.Server.Accept), zap.String("port", cfg.Server.Port), zap.Bool("https", cfg.Server.Secure))
	serverCfg := &server.Config{
		Production: cfg.Logger.Production,
		Logger:     logger,
		Accept:     cfg.Server.Accept,
		Port:       cfg.Server.Port,
	}

	messaging.Send(cfg.FreeClimb.AccountID, cfg.FreeClimb.AuthToken, cfg.FreeClimb.From, cfg.FreeClimb.To)

	s := StartServer(serverCfg, cfg.Server.Secure, cfg.Server.CertFile, cfg.Server.KeyFile)

	err = AwaitShutdown(logger, s)
	if err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	}
	logger.Info("shutdown")
}

func makeLogger(production bool, level string) (*zap.Logger, error) {
	var loggerConfig zap.Config
	if production {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	err := loggerConfig.Level.UnmarshalText([]byte(level))
	if err != nil {
		return nil, errstack.Push(err, "failed to unmarshal string version of logger level")
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, errstack.Push(err, "failed to build logger from generated config")
	}

	return logger, nil
}

func StartServer(serverCfg *server.Config, secure bool, certFile string, keyFile string) *server.Server {
	s := server.New(serverCfg)
	go s.Run(secure, certFile, keyFile)
	return s
}

func AwaitShutdown(logger *zap.Logger, s *server.Server) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	exitSignal := <-c

	logger.Info("received shutdown signal", zap.Any("exitSignal", exitSignal))

	return s.Shutdown(5 * time.Second)
}
