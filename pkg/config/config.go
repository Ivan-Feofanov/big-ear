package config

import (
	"braces.dev/errtrace"
	"github.com/cristalhq/aconfig"
)

type StreamConfig struct {
	MaxMsgs    int64 `default:"-1"`
	MaxBytes   int64 `default:"-1"`
	StreamName string
}

type Streams struct {
	Ethereum StreamConfig
	Polygon  StreamConfig
}

type Config struct {
	NatsURL           string `default:"nats://localhost:4222" env:"NATS_URL"`
	JsonRpcURL        string `required:"true" env:"JSON_RPC_URL"`
	PolygonJsonRpcURL string `required:"true" env:"POLYGON_JSON_RPC_URL"`
	Streams           Streams
	LastDataDir       string `default:"/tmp/"`
	Block             uint64
	Debug             bool
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

	var ethCfg StreamConfig
	ethCfgLoader := aconfig.LoaderFor(&ethCfg, aconfig.Config{
		SkipFlags:        true,
		EnvPrefix:        "ETH",
		AllowUnknownEnvs: true,
	})
	if err := ethCfgLoader.Load(); err != nil {
		return nil, errtrace.Wrap(err)
	}
	cfg.Streams.Ethereum = ethCfg

	var polygonCfg StreamConfig
	polygonCfgLoader := aconfig.LoaderFor(&polygonCfg, aconfig.Config{
		SkipFlags:        true,
		EnvPrefix:        "POLYGON",
		AllowUnknownEnvs: true,
	})
	if err := polygonCfgLoader.Load(); err != nil {
		return nil, errtrace.Wrap(err)
	}
	cfg.Streams.Polygon = polygonCfg
	return &cfg, nil
}
