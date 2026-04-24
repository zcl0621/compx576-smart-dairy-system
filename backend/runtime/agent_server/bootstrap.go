package agent_server

import (
	"context"
	"fmt"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server/consumer"
	alertservice "github.com/zcl0621/compx576-smart-dairy-system/service/alert"
	"go.uber.org/zap"
)

func Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", config.Get().AgentServer.Port)

	go consumer.StartMetricWriter(ctx)
	go consumer.StartAlertEvaluator(ctx)
	go alertservice.StartOfflineTicker(ctx)

	r := NewRouter()
	projectlog.L().Info("start agent server", zap.String("addr", addr))
	return r.Run(addr)
}
