package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/tendermint/tendermint/proto/tendermint/types"
	"google.golang.org/grpc"
)

// GetNodeInfo queries the current node info.
func (c *Client) GetNodeInfo(ctx context.Context, in *tmservice.GetNodeInfoRequest, opts ...grpc.CallOption) (*tmservice.GetNodeInfoResponse, error) {
	panic("not implemented")
}

// GetSyncing queries node syncing.
func (c *Client) GetSyncing(ctx context.Context, in *tmservice.GetSyncingRequest, opts ...grpc.CallOption) (*tmservice.GetSyncingResponse, error) {
	panic("not implemented")
}

// GetLatestBlock returns the latest block.
func (c *Client) GetLatestBlock(ctx context.Context, in *tmservice.GetLatestBlockRequest, opts ...grpc.CallOption) (*tmservice.GetLatestBlockResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	reps := tmservice.GetLatestBlockResponse{
		Block: &types.Block{
			Header: c.ctx.BlockHeader(), // We only need the header for simulation.
		},
	}

	return &reps, nil
}

// GetBlockByHeight queries block for given height.
func (c *Client) GetBlockByHeight(ctx context.Context, in *tmservice.GetBlockByHeightRequest, opts ...grpc.CallOption) (*tmservice.GetBlockByHeightResponse, error) {
	panic("not implemented")
}

// GetLatestValidatorSet queries latest validator-set.
func (c *Client) GetLatestValidatorSet(ctx context.Context, in *tmservice.GetLatestValidatorSetRequest, opts ...grpc.CallOption) (*tmservice.GetLatestValidatorSetResponse, error) {
	panic("not implemented")
}

// GetValidatorSetByHeight queries validator-set at a given height.
func (c *Client) GetValidatorSetByHeight(ctx context.Context, in *tmservice.GetValidatorSetByHeightRequest, opts ...grpc.CallOption) (*tmservice.GetValidatorSetByHeightResponse, error) {
	panic("not implemented")
}
