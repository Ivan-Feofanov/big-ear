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

	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}

}
