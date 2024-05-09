package main

import (
	"log"

	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	service "github.com/Ivan-Feofanov/big-ear/pkg/svc"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	svc, err := service.NewService(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting stream with name: %s...\n", cfg.Nats.StreamName)
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}

}
