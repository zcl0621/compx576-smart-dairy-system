package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent"
	"go.uber.org/zap"
)

func main() {
	host := flag.String("host", "127.0.0.1", "agent server host")
	port := flag.Int("port", 8081, "agent server port")
	cowID := flag.String("cow", "", "cow_id to simulate (required)")
	status := flag.String("status", "health", "health status: health / unhealth / ill")
	flag.Parse()

	if *cowID == "" {
		fmt.Fprintln(os.Stderr, "--cow flag is required")
		os.Exit(1)
	}

	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := projectlog.Init(); err != nil {
		panic(err)
	}
	defer projectlog.Sync()

	baseURL := fmt.Sprintf("http://%s:%d", *host, *port)

	projectlog.L().Info("agent starting",
		zap.String("cow_id", *cowID),
		zap.String("status", *status),
		zap.String("server", baseURL),
	)

	// get JWT token from agent server
	token, err := agent.FetchToken(baseURL, *cowID)
	if err != nil {
		projectlog.L().Fatal("failed to get token", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.StartHealthAgent(ctx, *cowID, baseURL, token, *status)
	go agent.StartMilkingAgent(ctx, *cowID, baseURL, token)
	go agent.StartWeightAgent(ctx, *cowID, baseURL, token)

	projectlog.L().Info("all agents running", zap.String("cow_id", *cowID))

	// wait for signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	projectlog.L().Info("shutting down agents")
	cancel()
}
