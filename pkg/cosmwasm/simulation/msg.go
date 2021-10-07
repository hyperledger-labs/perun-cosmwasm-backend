package simulation

import (
	"context"
	"fmt"

	wtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"google.golang.org/grpc"
)

// StoreCode stores the code of a smart contract on the ledger.
func (c *Client) StoreCode(ctx context.Context, in *wtypes.MsgStoreCode, opts ...grpc.CallOption) (*wtypes.MsgStoreCodeResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_ctx := c.ctx.WithContext(ctx)
	res, err := c.msgHandler(_ctx, in)
	if err != nil {
		return nil, fmt.Errorf("handling message: %w", err)
	}

	var resp wtypes.MsgStoreCodeResponse
	err = resp.Unmarshal(res.Data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return &resp, nil
}

//  InstantiateContract creates a new smart contract instance for the given code id.
func (c *Client) InstantiateContract(ctx context.Context, in *wtypes.MsgInstantiateContract, opts ...grpc.CallOption) (*wtypes.MsgInstantiateContractResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_ctx := c.ctx.WithContext(ctx)
	res, err := c.msgHandler(_ctx, in)
	if err != nil {
		return nil, fmt.Errorf("handling message: %w", err)
	}

	var resp wtypes.MsgInstantiateContractResponse
	err = resp.Unmarshal(res.Data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return &resp, nil
}

// ExecuteContract executes a function on a contract.
func (c *Client) ExecuteContract(ctx context.Context, in *wtypes.MsgExecuteContract, opts ...grpc.CallOption) (*wtypes.MsgExecuteContractResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_ctx := c.ctx.WithContext(ctx)
	res, err := c.msgHandler(_ctx, in)
	if err != nil {
		return nil, fmt.Errorf("handling message: %w", err)
	}

	var resp wtypes.MsgExecuteContractResponse
	err = resp.Unmarshal(res.Data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return &resp, nil
}

// MigrateContract performs a code upgrade or downgrade for a smart contract.
func (c *Client) MigrateContract(ctx context.Context, in *wtypes.MsgMigrateContract, opts ...grpc.CallOption) (*wtypes.MsgMigrateContractResponse, error) {
	panic("not implemented")
}

// UpdateAdmin sets a new admin for a smart contract.
func (c *Client) UpdateAdmin(ctx context.Context, in *wtypes.MsgUpdateAdmin, opts ...grpc.CallOption) (*wtypes.MsgUpdateAdminResponse, error) {
	panic("not implemented")
}

// ClearAdmin removes any admin stored for a smart contract.
func (c *Client) ClearAdmin(ctx context.Context, in *wtypes.MsgClearAdmin, opts ...grpc.CallOption) (*wtypes.MsgClearAdminResponse, error) {
	panic("not implemented")
}
