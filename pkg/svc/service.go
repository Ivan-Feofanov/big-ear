package service

import (
	"github.com/Ivan-Feofanov/go-blockchain-listener/pkg/config"

	protocol "github.com/Ivan-Feofanov/go-blockchain-listener/pkg/svc/.generated"

	"google.golang.org/grpc/credentials/insecure"

	"braces.dev/errtrace"
	"google.golang.org/grpc"
)

type Puller interface {
	Pull(output protocol.AgentClient) error
}

type Service struct {
	eth *EthClient
	out protocol.AgentClient
}

func NewService(cfg *config.Config) (*Service, error) {
	eth, err := NewEthClient(cfg)
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
