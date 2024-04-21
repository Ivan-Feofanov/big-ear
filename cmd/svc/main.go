package main

import (
	"log/slog"
	"os"

	"github.com/Ivan-Feofanov/go-blockchain-listener/pkg/config"
	service "github.com/Ivan-Feofanov/go-blockchain-listener/pkg/svc"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	svc, err := service.NewService(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := svc.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

}
