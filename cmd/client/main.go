package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"sinergy-test/cmd/config"
	"sinergy-test/internal/client/service"
	"sinergy-test/internal/client/transport/http"
)

func init() {
	time.Local = time.UTC
}

func main() {
	cfg := config.GetConfig("dev.env")

	fmt.Println("serverHost:", cfg.ServerHost)
	fmt.Println("serverPort:", cfg.ServerPort)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	srv := service.NewService(sugar, cfg.ServerHost, cfg.ServerPort)
	ctrl := http.NewController(sugar, srv)
	err := ctrl.Run(cfg.ClientPort)
	if err != nil {
		panic(err)
	}
}
