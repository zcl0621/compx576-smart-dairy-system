// @title smart dairy api
// @version 1.0
// @description api for smart dairy system
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	stdlog "log"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	webserver "github.com/zcl0621/compx576-smart-dairy-system/runtime/web_server"
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

	if err := webserver.Run(); err != nil {
		stdlog.Fatal(err)
	}
}
