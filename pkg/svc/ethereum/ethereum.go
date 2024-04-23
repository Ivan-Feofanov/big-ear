package ethereum

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	protocol "github.com/Ivan-Feofanov/big-ear/pkg/svc/protocol"

	"braces.dev/errtrace"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	cfg      *config.Config
	URL      string
	client   EthClient
	runLimit uint
}

func GetClient(cfg *config.Config, client EthClient, dialURL string, runLimit uint) *Client {
	return &Client{
		cfg:      cfg,
		URL:      dialURL,
		client:   client,
		runLimit: runLimit,
	}
}

func New(cfg *config.Config) (*Client, error) {
	client, err := ethclient.Dial(buildEthURL(cfg))
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return GetClient(cfg, client, buildEthURL(cfg), 0), nil

}

func (e *Client) Pull(out protocol.AgentClient) error {
	var requestID uint = 0
	for {
		requestID++

		blockNumber, err := e.client.BlockNumber(context.Background())
		if err != nil {
			log.Default().Println("Failed to retrieve block number:", err)
			return errtrace.Wrap(err)
		}
		block, err := e.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		if err != nil {
			return errtrace.Wrap(err)
		}
		blockEvent := protocol.BlockEvent{
			BlockNumber: block.Number().String(),
			Block: &protocol.BlockEvent_EthBlock{
				Hash:       block.Hash().String(),
				ParentHash: block.ParentHash().String(),
				Number:     block.Number().String(),
			},
		}
		resp, err := out.EvaluateBlock(context.Background(), &protocol.EvaluateBlockRequest{
			Event:     &blockEvent,
			RequestId: strconv.Itoa(int(requestID)),
			ShardId:   1,
		})
		if err != nil {
			return errtrace.Wrap(err)
		}

		if resp.GetStatus() != protocol.ResponseStatus_SUCCESS {
			log.Default().Printf("failed to evaluate block #%d", blockNumber)
			var errs []error
			for _, e := range resp.GetErrors() {
				errs = append(errs, errors.New(e.Message))
			}
			return errtrace.Wrap(errors.Join(errs...))
		}

		if e.cfg.Verbosity == config.VerbosityDebug {
			log.Default().Println(resp.Status)
			log.Default().Println(resp.GetErrors())
			log.Default().Println(resp.GetFindings())
		}

		if e.runLimit > 0 && requestID >= e.runLimit {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func buildEthURL(cfg *config.Config) string {
	return fmt.Sprintf("https://lb.drpc.org/ogrpc?network=ethereum&dkey=%s", cfg.DRPCAPIKey)
}
