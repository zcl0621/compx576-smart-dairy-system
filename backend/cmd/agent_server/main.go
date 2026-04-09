package main

import (
	"context"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	agent_server "github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := projectlog.Init(); err != nil {
		panic(err)
	}
	defer projectlog.Sync()

	if err := pg.InitDB(); err != nil {
		panic(err)
	}
	if err := model.Migrate(pg.DB); err != nil {
		panic(err)
	}
	if err := redis.InitRedis(); err != nil {
		panic(err)
	}
	if err := mq.Init(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err := agent_server.Run(ctx); err != nil {
		panic(err)
	}
}
