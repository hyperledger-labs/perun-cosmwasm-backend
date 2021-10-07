//  Copyright 2021 PolyCrypt GmbH
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package simulation

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// Client is a Cosmos client that can be used in a simulation environment.
type Client struct {
	chainAccount types.AccAddress
	msgHandler   types.Handler
	ctx          types.Context
	keepers      keeper.TestKeepers
	mu           sync.RWMutex
	stopTick     chan struct{}
}

var _ client.Client = &Client{}

// NewTestClient creates a new client specifically for a test environment.
func NewTestClient(t *testing.T) *Client {
	ctx, keepers := wasm.CreateTestInput(t, false, "staking,stargate")
	cdc := keeper.MakeTestCodec(t)
	mod := wasm.NewAppModule(cdc, keepers.WasmKeeper, keepers.StakingKeeper)
	handler := mod.Route().Handler()
	acc := createAccount(ctx, keepers.AccountKeeper)

	return NewClient(ctx, acc, handler, keepers)
}

func NewClient(ctx types.Context, acc types.AccAddress, handler types.Handler, keepers keeper.TestKeepers) *Client {
	return &Client{
		chainAccount: acc,
		msgHandler:   handler,
		ctx:          ctx,
		keepers:      keepers,
	}
}

// Account returns the account that is used for sending transactions.
func (c *Client) Account() types.AccAddress {
	return c.chainAccount
}

func createAccount(ctx types.Context, am authkeeper.AccountKeeper) types.AccAddress {
	_, _, addr := keyPubAddr()
	acc := am.NewAccountWithAddress(ctx, addr)
	am.SetAccount(ctx, acc)
	return addr
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, types.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := types.AccAddress(pub.Address())
	return key, pub, addr
}

// AddCoins adds the specified amount of coins to the specified account.
func (c *Client) AddCoins(ctx context.Context, addr types.AccAddress, coins types.Coins) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_ctx := c.ctx.WithContext(ctx)
	return c.keepers.BankKeeper.AddCoins(_ctx, addr, coins)
}

// BlockTime returns the block time.
func (c *Client) BlockTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.ctx.BlockTime()
}

// SetBlockTime sets the block time.
func (c *Client) SetBlockTime(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx = c.ctx.WithBlockTime(t).WithBlockHeight(c.ctx.BlockHeight() + 1)
}

// StartTicking starts the auto-ticking process. The process will create a
// new block and advance the time on every tick.
func (c *Client) StartTicking(tickInterval time.Duration, timeAddedPerTick time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopTick != nil {
		panic("already ticking")
	}

	c.stopTick = make(chan struct{})
	go func() {
		for {
			select {
			case <-c.stopTick:
				c.mu.Lock()
				c.stopTick = nil
				c.mu.Unlock()
				return
			case <-time.After(tickInterval):
				func() {
					c.mu.Lock()
					defer c.mu.Unlock()
					t := c.ctx.BlockTime().Add(timeAddedPerTick)
					h := c.ctx.BlockHeight() + 1
					c.ctx = c.ctx.WithBlockTime(t).WithBlockHeight(h)
				}()
			}
		}
	}()
}

// StopTicking stops the auto-ticking process.
func (c *Client) StopTicking() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	close(c.stopTick)
}

// StoreContractTemplate stores a contract on the blockchain.
func (c *Client) StoreContractTemplate(ctx context.Context, contract client.ContractTemplate) (client.StoredContract, error) {
	msg := &wtypes.MsgStoreCode{
		Sender:       c.chainAccount.String(),
		WASMByteCode: contract.Code(),
	}

	resp, err := c.StoreCode(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("storing code: %w", err)
	}

	return client.NewStoredContract(contract, resp.CodeID), nil
}

// InstantiateStoredContract creates a contract instance from a stored contract.
func (c *Client) InstantiateStoredContract(ctx context.Context, contract client.StoredContract, msg []byte, deposit types.Coins) (client.ContractInstance, []byte, error) {
	err := contract.ValidateInitMsg(msg)
	if err != nil {
		return nil, nil, err
	}

	_msg := &wtypes.MsgInstantiateContract{
		Sender: c.chainAccount.String(),
		CodeID: contract.ID(),
		Msg:    msg,
		Funds:  deposit,
	}

	resp, err := c.InstantiateContract(ctx, _msg)
	if err != nil {
		return nil, nil, err
	}

	return client.NewContractInstance(contract, resp.Address), resp.Data, nil
}
