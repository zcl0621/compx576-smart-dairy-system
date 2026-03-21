package main

import (
	stdlog "log"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	if err := projectlog.Init(); err != nil {
		panic(err)
	}

	defer projectlog.Sync()

	stdlog.Println("agent server stub is up")
}
