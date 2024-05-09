package service

import (
	"braces.dev/errtrace"
	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/Ivan-Feofanov/big-ear/pkg/ethereum"
	"github.com/Ivan-Feofanov/big-ear/pkg/stream"
	"github.com/nats-io/nats.go"
	"github.com/sourcegraph/conc/pool"
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

	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	eventStream, err := stream.NewStream(nc, cfg.Nats)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &Service{
		eth:    eth,
		stream: eventStream,
	}, nil
}

func (s *Service) Run() error {
	pullers := []Puller{s.eth}

	p := pool.New().WithErrors()
	for _, puller := range pullers {
		p.Go(func() error {
			return puller.Pull(s.stream)
		})
	}
	return errtrace.Wrap(p.Wait())
}
