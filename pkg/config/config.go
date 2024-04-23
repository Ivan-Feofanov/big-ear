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
	Verbosity       int    `default:"0" usage:"verbosity level"`
	DRPCAPIKey      string `required:"true" env:"DRPC_API_KEY" toml:"drpc_api_key"`
	ReceiverAddress string `default:"localhost:50051"`
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
