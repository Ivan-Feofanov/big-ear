package config

import (
	"braces.dev/errtrace"
	"github.com/cristalhq/aconfig"
)

type NatsConfig struct {
	URL        string `default:"nats://localhost:4222"`
	MaxMsgs    int64  `default:"-1"`
	MaxBytes   int64  `default:"-1"`
	StreamName string `required:"true"`
}

type Config struct {
	JsonRpcURL  string `required:"true" env:"JSON_RPC_URL"`
	Nats        NatsConfig
	LastDataDir string `default:"/tmp/"`
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
