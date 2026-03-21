package web_server

import (
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"go.uber.org/zap"
)

const defaultAddr = ":8080"

func Run() error {

	r := NewRouter()
	projectlog.L().Info("start web", zap.String("addr", defaultAddr))

	return r.Run(defaultAddr)
}
