package polygon

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"braces.dev/errtrace"
	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/Ivan-Feofanov/big-ear/pkg/ethereum"
	protocol "github.com/Ivan-Feofanov/big-ear/pkg/proto"
	"github.com/Ivan-Feofanov/big-ear/pkg/stream"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/protobuf/proto"
)

const lastBlockFilename = "last_block_polygon.txt"

type Client struct {
	cfg             *config.Config
	client          ethereum.EthClient
	lastBlockNumber uint64
}

func GetClient(cfg *config.Config, client ethereum.EthClient) *Client {
	return &Client{
		cfg:    cfg,
		client: client,
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

	client, err := ethclient.Dial(cfg.PolygonJsonRpcURL)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return GetClient(cfg, client), nil
}

func (e *Client) Pull(eventStream *stream.Stream) error {
	if e.cfg.Block != 0 {
		return e.PullBlock(e.cfg.Block, eventStream)
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
	}
}

func (e *Client) StreamBlock(block types.Block, eventStream *stream.Stream) error {
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
	if e.cfg.Debug {
		log.Default().Printf("Block #%d streamed", block.Number())
	}
	return nil
}

func (e *Client) StreamTransaction(tx *types.Transaction, eventStream *stream.Stream) error {
	receipt, err := e.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return errtrace.Wrap(err)
	}

	txEvent := protocol.TransactionEvent{
		Block: &protocol.TransactionEvent_EthBlock{
			BlockNumber: receipt.BlockNumber.String(),
			BlockHash:   receipt.BlockHash.String(),
		},
		Addresses: map[string]bool{strings.ToLower(from.String()): true},
		//Logs:      convertLogs(receipt.Logs),
	}
	for _, record := range receipt.Logs {
		txEvent.Addresses[strings.ToLower(record.Address.String())] = true
	}

	payload, err := proto.Marshal(&protocol.EvaluateTxRequest{
		Event:   &txEvent,
		ShardId: 1,
	})
	if err != nil {
		return errtrace.Wrap(err)
	}

	if err = eventStream.Publish(context.Background(), "transaction", payload); err != nil {
		return errtrace.Wrap(err)
	}
	if e.cfg.Debug {
		log.Default().Printf("Tx #%s streamed", tx.Hash().String())
	}
	return nil
}

func (e *Client) PullBlock(blockNumber uint64, eventStream *stream.Stream) error {
	block, err := e.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return errtrace.Wrap(err)
	}
	transactions := block.Transactions()
	for _, tx := range transactions {
		if err := e.StreamTransaction(tx, eventStream); err != nil {
			return errtrace.Wrap(err)
		}
	}
	return e.StreamBlock(*block, eventStream)
}
