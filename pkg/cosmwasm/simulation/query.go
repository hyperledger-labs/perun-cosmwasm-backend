package simulation

import (
	"context"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
)

// ContractInfo gets the contract meta data.
func (c *Client) ContractInfo(ctx context.Context, in *wtypes.QueryContractInfoRequest, opts ...grpc.CallOption) (*wtypes.QueryContractInfoResponse, error) {
	panic("not implemented")
}

// ContractHistory gets the contract code history.
func (c *Client) ContractHistory(ctx context.Context, in *wtypes.QueryContractHistoryRequest, opts ...grpc.CallOption) (*wtypes.QueryContractHistoryResponse, error) {
	panic("not implemented")
}

// ContractsByCode lists all smart contracts for a code id.
func (c *Client) ContractsByCode(ctx context.Context, in *wtypes.QueryContractsByCodeRequest, opts ...grpc.CallOption) (*wtypes.QueryContractsByCodeResponse, error) {
	panic("not implemented")
}

// AllContractState gets all raw store data for a single contract.
func (c *Client) AllContractState(ctx context.Context, in *wtypes.QueryAllContractStateRequest, opts ...grpc.CallOption) (*wtypes.QueryAllContractStateResponse, error) {
	panic("not implemented")
}

// RawContractState gets a single key from the raw store data of a contract.
func (c *Client) RawContractState(ctx context.Context, in *wtypes.QueryRawContractStateRequest, opts ...grpc.CallOption) (*wtypes.QueryRawContractStateResponse, error) {
	panic("not implemented")
}

// SmartContractState performs a smart contract query.
func (c *Client) SmartContractState(ctx context.Context, in *wtypes.QuerySmartContractStateRequest, opts ...grpc.CallOption) (*wtypes.QuerySmartContractStateResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	q := keeper.Querier(c.keepers.WasmKeeper)
	_ctx := c.ctx.WithContext(ctx)
	return q.SmartContractState(
		types.WrapSDKContext(_ctx),
		in,
	)
}

// Code gets the binary code and metadata for a code id.
func (c *Client) Code(ctx context.Context, in *wtypes.QueryCodeRequest, opts ...grpc.CallOption) (*wtypes.QueryCodeResponse, error) {
	panic("not implemented")
}

// Codes gets the metadata for all stored wasm codes.
func (c *Client) Codes(ctx context.Context, in *wtypes.QueryCodesRequest, opts ...grpc.CallOption) (*wtypes.QueryCodesResponse, error) {
	panic("not implemented")
}
