package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	protocol "github.com/Ivan-Feofanov/big-ear/pkg/proto"
	"github.com/Ivan-Feofanov/big-ear/pkg/stream"
	"google.golang.org/protobuf/proto"

	"braces.dev/errtrace"
	"github.com/ethereum/go-ethereum/ethclient"
)

const lastBlockFilename = "last_block.txt"

type Client struct {
	cfg             *config.Config
	URL             string
	client          EthClient
	lastBlockNumber uint64
	runLimit        uint
}

func GetClient(cfg *config.Config, client EthClient, dialURL string, runLimit uint) *Client {
	return &Client{
		cfg:      cfg,
		URL:      dialURL,
		client:   client,
		runLimit: runLimit,
	}
}

func checkDataDir(dataDir string) error {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return errtrace.Wrap(err)
		}
	}
	return nil
}

func New(cfg *config.Config) (*Client, error) {
	if err := checkDataDir(cfg.LastDataDir); err != nil {
		return nil, errtrace.Wrap(err)
	}

	client, err := ethclient.Dial(buildEthURL(cfg))
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return GetClient(cfg, client, buildEthURL(cfg), 0), nil

}

func (e *Client) PullBlock(blockNumber uint64, eventStream *stream.Stream) error {
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

	payload, err := proto.Marshal(&protocol.EvaluateBlockRequest{
		Event:   &blockEvent,
		ShardId: 1,
	})
	if err != nil {
		return errtrace.Wrap(err)
	}

	if err = eventStream.Publish(context.Background(), "block", payload); err != nil {
		return errtrace.Wrap(err)
	}

	log.Default().Printf("Block #%d streamed", block.Number())

	return nil
}

func readLastBlockNumber(dataDir string) (uint64, error) {
	data, err := os.ReadFile(path.Join(dataDir, lastBlockFilename))
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	num, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	return num, nil
}

func writeLastBlockNumber(dataDir string, num uint64) error {
	return os.WriteFile(path.Join(dataDir, lastBlockFilename), []byte(strconv.FormatUint(num, 10)), 0644)
}

func (e *Client) Rewind(from, to uint64, eventStream *stream.Stream) error {
	log.Default().Printf("Service was off since block #%d. Rewinding to actual #%d", from, to)
	for i := from + 1; i <= to; i++ {
		if err := e.PullBlock(i, eventStream); err != nil {
			return errtrace.Wrap(err)
		}
	}
	if err := writeLastBlockNumber(e.cfg.LastDataDir, to); err != nil {
		log.Default().Printf("Failed to write last block number to file: %v", err)
		return errtrace.Wrap(err)
	}

	return nil
}

func (e *Client) Pull(eventStream *stream.Stream) error {
	actualBlockNumber, err := e.client.BlockNumber(context.Background())
	if err != nil {
		return errtrace.Wrap(err)
	}
	writtenLastBlockNumber, err := readLastBlockNumber(e.cfg.LastDataDir)
	if err != nil {
		log.Default().Printf(
			"Failed to read last block number from file: %v\nStaring from actual block number #%d", err, actualBlockNumber)
		writtenLastBlockNumber = actualBlockNumber
	}

	if writtenLastBlockNumber < actualBlockNumber {
		if err = e.Rewind(writtenLastBlockNumber, actualBlockNumber, eventStream); err != nil {
			return errtrace.Wrap(err)
		}
		e.lastBlockNumber = actualBlockNumber
		return e.Pull(eventStream)
	}

	var requestID uint = 0
	for {
		requestID++

		blockNumber, err := e.client.BlockNumber(context.Background())
		if err != nil {
			log.Default().Println("Failed to retrieve block number:", err)
			return errtrace.Wrap(err)
		}
		if blockNumber == e.lastBlockNumber {
			time.Sleep(3 * time.Second)
			continue
		}
		e.lastBlockNumber = blockNumber

		if err = e.PullBlock(blockNumber, eventStream); err != nil {
			return errtrace.Wrap(err)
		}

		if err := writeLastBlockNumber(e.cfg.LastDataDir, blockNumber); err != nil {
			log.Default().Printf("Failed to write last block number to file: %v", err)
			return errtrace.Wrap(err)
		}

		if e.runLimit > 0 && requestID >= e.runLimit {
			break
		}
	}

	return nil
}

func buildEthURL(cfg *config.Config) string {
	return fmt.Sprintf("https://lb.drpc.org/ogrpc?network=ethereum&dkey=%s", cfg.DRPCAPIKey)
}
