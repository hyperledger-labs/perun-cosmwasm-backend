package channel

import (
	"context"

	wtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
)

type contractClient struct {
	client   client.Client
	contract client.ContractInstance
	acc      types.AccAddress
}

func newContractClient(client client.Client, contract client.ContractInstance, acc types.AccAddress) *contractClient {
	return &contractClient{
		client:   client,
		contract: contract,
		acc:      acc,
	}
}

// Account returns the account used for interacting with the network.
func (c *contractClient) Account() types.AccAddress {
	return c.acc
}

// Query submits a contract query.
func (c *contractClient) Query(ctx context.Context, msg []byte) (*wtypes.QuerySmartContractStateResponse, error) {
	err := c.contract.ValidateQueryMsg(msg)
	if err != nil {
		return nil, err
	}

	req := &wtypes.QuerySmartContractStateRequest{
		Address:   c.contract.Address(),
		QueryData: msg,
	}
	return c.client.SmartContractState(ctx, req)
}

// Execute executes a contract function.
func (c *contractClient) Execute(ctx context.Context, msg []byte, funds types.Coins) (*wtypes.MsgExecuteContractResponse, error) {
	err := c.contract.ValidateExecuteMsg(msg)
	if err != nil {
		return nil, err
	}

	_msg := &wtypes.MsgExecuteContract{
		Sender:   c.acc.String(),
		Contract: c.contract.Address(),
		Msg:      msg,
		Funds:    funds,
	}
	return c.client.ExecuteContract(ctx, _msg)
}
