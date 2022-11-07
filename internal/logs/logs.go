package logs

import (
	"log"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/env"
	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger логгер необходимо инициализировать при запуске приложения
func InitLogger() {
	var err error
	if !env.InProd() {
		logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = cfg.Build()
	}
	if err != nil {
		log.Fatal("cannot init zap", err)
	}
}

func Error(text string, fields ...zap.Field) {
	logger.Error(text, fields...)
}

func Info(text string, fields ...zap.Field) {
	logger.Info(text, fields...)
}
