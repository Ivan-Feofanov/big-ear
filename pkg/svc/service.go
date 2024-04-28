package service

import (
	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/Ivan-Feofanov/big-ear/pkg/ethereum"
	"github.com/Ivan-Feofanov/big-ear/pkg/stream"
	"github.com/nats-io/nats.go"

	"braces.dev/errtrace"
)

type Puller interface {
	Pull(*stream.Stream) error
}

type Service struct {
	eth    *ethereum.Client
	stream *stream.Stream
}

func NewService(cfg *config.Config) (*Service, error) {
	eth, err := ethereum.New(cfg)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	eventStream, err := stream.NewStream(nc, cfg.StreamName)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &Service{
		eth:    eth,
		stream: eventStream,
	}, nil
}

func (s *Service) Run() error {
	return errtrace.Wrap(s.eth.Pull(s.stream))
}
