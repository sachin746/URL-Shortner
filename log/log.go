package log

import (
	"context"
	"log"

	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func InitLogger(ctx context.Context) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	Sugar = logger.Sugar()

	Sugar.Infof("Logger initialized")
}
