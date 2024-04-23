package service

import (
	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/Ivan-Feofanov/big-ear/pkg/svc/ethereum"

	protocol "github.com/Ivan-Feofanov/big-ear/pkg/svc/protocol"

	"google.golang.org/grpc/credentials/insecure"

	"braces.dev/errtrace"
	"google.golang.org/grpc"
)

type Puller interface {
	Pull(output protocol.AgentClient) error
}

type Service struct {
	eth *ethereum.Client
	out protocol.AgentClient
}

func NewService(cfg *config.Config) (*Service, error) {
	eth, err := ethereum.New(cfg)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	conn, err := grpc.Dial(cfg.ReceiverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &Service{
		eth: eth,
		out: protocol.NewAgentClient(conn),
	}, nil
}

func (s *Service) Run() error {
	return errtrace.Wrap(s.eth.Pull(s.out))
}
