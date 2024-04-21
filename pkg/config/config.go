package config

import (
	"braces.dev/errtrace"
	"github.com/cristalhq/aconfig"
)

const (
	VerbosityInfo  = 0
	VerbosityDebug = 1
)

type Config struct {
	Verbosity       int    `default:"1" usage:"verbosity level" env:"VERBOSITY" toml:"verbosity" flag:"v"`
	Environment     string `default:"dev" usage:"environment" env:"ENV"`
	Network         string `default:"mainnet" env:"NETWORK"`
	DRPCAPIKey      string `required:"true" env:"DRPC_API_KEY" toml:"drpc_api_key" flag:"drpc-api-key"`
	ReceiverAddress string `default:"localhost:50051" env:"RECEIVER_ADDRESS"`
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
