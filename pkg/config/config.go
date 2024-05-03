package config

import (
	"braces.dev/errtrace"
	"github.com/cristalhq/aconfig"
)

type Config struct {
	DRPCAPIKey  string `required:"true" env:"DRPC_API_KEY"`
	NatsURL     string `default:"nats://localhost:4222"`
	LastDataDir string `default:"/tmp/"`
	StreamName  string `required:"true" env:"STREAM_NAME"`
	Block       uint64
	Debug       bool
}

func GetConfig() (*Config, error) {
	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		AllowUnknownFlags: true,
		SkipFlags:         true,
	})

	if err := loader.Load(); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &cfg, nil
}
