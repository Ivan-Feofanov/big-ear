package service

import (
	"log"

	"braces.dev/errtrace"
	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/Ivan-Feofanov/big-ear/pkg/ethereum"
	"github.com/Ivan-Feofanov/big-ear/pkg/polygon"
	"github.com/Ivan-Feofanov/big-ear/pkg/stream"
	"github.com/nats-io/nats.go"
	"github.com/sourcegraph/conc/pool"
)

type Puller interface {
	Pull(*stream.Stream) error
}

type Service struct {
	eth     *ethereum.Client
	polygon *polygon.Client
	stream  *stream.Stream
}

func NewService(cfg *config.Config) (*Service, error) {
	eth, err := ethereum.New(cfg)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	poly, err := polygon.New(cfg)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	eventStream, err := stream.NewStream(nc, cfg.Streams.Polygon)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &Service{
		eth:     eth,
		polygon: poly,
		stream:  eventStream,
	}, nil
}

func (s *Service) Run() error {
	pullers := []Puller{s.polygon}

	p := pool.New().WithErrors()
	for _, puller := range pullers {
		p.Go(func() error {
			log.Printf("Starting stream with name: %s...\n", s.stream.Name())
			return puller.Pull(s.stream)
		})
	}
	return errtrace.Wrap(p.Wait())
}
