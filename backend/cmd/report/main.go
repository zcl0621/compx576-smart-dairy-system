package main

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	redisdb "github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/service/llm"
	"github.com/zcl0621/compx576-smart-dairy-system/service/report"
)

func main() {
	if err := config.InitConfig(); err != nil {
		stdlog.Printf("config: %v", err)
		os.Exit(1)
	}
	if err := projectlog.Init(); err != nil {
		stdlog.Printf("log: %v", err)
		os.Exit(1)
	}
	defer projectlog.Sync()

	if err := pg.InitDB(); err != nil {
		projectlog.L().Error("pg init", zap.Error(err))
		os.Exit(1)
	}
	if err := redisdb.InitRedis(); err != nil {
		projectlog.L().Error("redis init", zap.Error(err))
		os.Exit(1)
	}

	cfg := config.Get().Deepseek
	client := &llm.Client{
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
		Timeout: cfg.Timeout,
		HTTP:    &http.Client{Timeout: cfg.Timeout},
	}
	gen := &report.Generator{
		DB:  pg.DB,
		LLM: client,
		Now: time.Now,
	}

	ctx := context.Background()
	if err := gen.RunOnce(ctx); err != nil {
		projectlog.L().Error("report run failed", zap.Error(err))
		os.Exit(1)
	}
	projectlog.L().Info("report run finished")
}
