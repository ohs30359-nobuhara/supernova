package pkg

import (
	"go.uber.org/zap"
	"sync"
)

var once sync.Once
var logger *zap.Logger

func GetLogger() *zap.Logger {
	once.Do(func() {
		// Loggerの初期化
		cfg := zap.NewDevelopmentConfig()

		var err error
		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
	})

	return logger
}
