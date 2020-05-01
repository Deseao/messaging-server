package messaging

import (
	"time"

	"go.uber.org/zap"

	"github.com/Deseao/messaging-server/internal/config"
	"github.com/Deseao/messaging-server/internal/state"
)

func RunQueue(logger *zap.Logger, info *state.State, cfg config.FreeClimb) {
	for now := range time.Tick(3 * time.Second) {
		logger.Info("Checking for unsent messages")
		for key, subscriber := range info.Subscribers {
			if now.After(subscriber.LastSent.Add(30 * time.Second)) {
				logger.Debug("Sending message", zap.String("number", subscriber.Number))
				subscriber.LastSent = now
				info.Subscribers[key] = subscriber
				Send(cfg.AccountID, cfg.AuthToken, cfg.From, subscriber.Number)
			}
		}
	}
}
