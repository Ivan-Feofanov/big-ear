package ethereum_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/Ivan-Feofanov/big-ear/pkg/svc/ethereum"
	ethereummock "github.com/Ivan-Feofanov/big-ear/pkg/svc/ethereum/mocks"
	"github.com/Ivan-Feofanov/big-ear/pkg/svc/protocol"
	protocolmock "github.com/Ivan-Feofanov/big-ear/pkg/svc/protocol/mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Pull(t *testing.T) {
	blockHeader := &types.Header{
		Number: big.NewInt(int64(1)),
	}
	block := types.NewBlockWithHeader(blockHeader)

	type errs struct {
		blockNumber error
		block       error
		evaluate    error
		gRPC        []*protocol.Error
	}
	tests := []struct {
		name       string
		errs       errs
		grpcStatus protocol.ResponseStatus
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "success",
			grpcStatus: protocol.ResponseStatus_SUCCESS,
			wantErr:    assert.NoError,
		},
		{
			name: "failed to retrieve block number",
			errs: errs{
				blockNumber: errors.New("failed to retrieve block number"),
			},
			wantErr: assert.Error,
		}, {
			name: "failed to retrieve block",
			errs: errs{
				block: errors.New("failed to retrieve block"),
			},
			wantErr: assert.Error,
		}, {
			name: "failed to evaluate block - general error",
			errs: errs{
				evaluate: errors.New("failed to evaluate block"),
			},
			wantErr: assert.Error,
		}, {
			name:       "failed to evaluate block - error response",
			grpcStatus: protocol.ResponseStatus_ERROR,
			errs: errs{
				gRPC: []*protocol.Error{{Message: "somethig went wrong"}},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := ethereummock.NewMockEthClient(t)
			clientMock.On("BlockNumber", mock.Anything).Return(uint64(1), tt.errs.blockNumber)
			clientMock.On("BlockByNumber", mock.Anything, big.NewInt(int64(1))).Return(block, tt.errs.block).Maybe()
			outMock := protocolmock.NewMockAgentClient(t)
			outMock.On("EvaluateBlock", mock.Anything, mock.Anything, mock.Anything).Return(&protocol.EvaluateBlockResponse{
				Status:    tt.grpcStatus,
				Errors:    tt.errs.gRPC,
				Findings:  nil,
				Metadata:  nil,
				Timestamp: "",
				LatencyMs: 0,
				Private:   false,
			}, tt.errs.evaluate).Maybe()
			e := ethereum.GetClient(&config.Config{}, clientMock, "", 1)

			tt.wantErr(t, e.Pull(outMock))
		})
	}
}
