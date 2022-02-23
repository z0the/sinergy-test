package main

import (
	"time"

	"go.uber.org/zap"

	"sinergy-test/internal/server/service"
	"sinergy-test/internal/server/transport/tcp"
)

func init() {
	time.Local = time.UTC
}

const (
	maxSleepTimeSeconds = 3
	serverPort          = "5001"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	resourceURLs := []string{
		"https://novasite.su/test1.php",
		"https://novasite.su/test2.php",
	}

	srv := service.NewService(sugar, resourceURLs, maxSleepTimeSeconds)

	ctrl := tcp.NewController(sugar, srv)

	err := ctrl.Run(serverPort)
	if err != nil {
		panic(err)
	}
}
