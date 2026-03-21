package log

import (
	"sync"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func Init() error {
	var initErr error
	once.Do(func() {
		cfg := zap.NewDevelopmentConfig()
		if config.Get().App.Env == "prod" {
			cfg = zap.NewProductionConfig()
		}

		built, err := cfg.Build()
		if err != nil {
			initErr = err
			return
		}

		logger = built
		zap.ReplaceGlobals(logger)
	})

	return initErr
}

func L() *zap.Logger {
	if logger == nil {
		return zap.L()
	}
	return logger
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}

func Sync() error {
	if logger == nil {
		return nil
	}
	return logger.Sync()
}
