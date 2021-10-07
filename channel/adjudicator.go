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

package channel

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

const defaultPollingInterval = 1 * time.Second

// Adjudicator provides methods for dispute resolution on the ledger.
type Adjudicator struct {
	*contractClient
	contract client.ContractInstance
	polling  time.Duration
}

type AdjudicatorOpt func(*Adjudicator)

func AdjudicatorPollingIntervalOpt(d time.Duration) AdjudicatorOpt {
	return func(a *Adjudicator) {
		a.polling = d
	}
}

func NewAdjudicator(c client.Client, contract client.ContractInstance, acc types.AccAddress, opts ...AdjudicatorOpt) *Adjudicator {
	a := &Adjudicator{
		contract:       contract,
		polling:        defaultPollingInterval,
		contractClient: newContractClient(c, contract, acc),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Register registers the given ledger channel state on-chain.
// If the channel has locked funds into sub-channels, the corresponding
// signed sub-channel states must be provided.
func (a *Adjudicator) Register(ctx context.Context, req channel.AdjudicatorReq, subChannels []channel.SignedState) error {
	if len(subChannels) > 0 {
		return fmt.Errorf("subchannels not supported")
	}
	return a.dispute(ctx, req)
}

func (a *Adjudicator) dispute(ctx context.Context, req channel.AdjudicatorReq) error {
	return a.callAdjudicator(ctx, binding.NewDisputeExecuteMsg, req)
}

type AdjudicatorMsgFunc func(p channel.Params, s channel.State, sigs []wallet.Sig) ([]byte, error)

func (a *Adjudicator) callAdjudicator(ctx context.Context, fn AdjudicatorMsgFunc, req channel.AdjudicatorReq) error {
	msg, err := fn(*req.Params, *req.Tx.State, req.Tx.Sigs)
	if err != nil {
		return err
	}

	_, err = a.Execute(ctx, msg, nil)
	return err
}

// Withdraw concludes and withdraws the registered state, so that the
// final outcome is set on the asset holders and funds are withdrawn.
// If the channel has locked funds in sub-channels, the states of the
// corresponding sub-channels need to be supplied additionally.
func (a *Adjudicator) Withdraw(ctx context.Context, req channel.AdjudicatorReq, subStates channel.StateMap) error {
	if len(subStates) > 0 {
		return fmt.Errorf("subchannels not supported")
	}
	err := a.conclude(ctx, req)
	if err != nil {
		return fmt.Errorf("concluding: %w", err)
	}
	return a.withdraw(ctx, req)
}

func (a *Adjudicator) conclude(ctx context.Context, req channel.AdjudicatorReq) error {
	return a.callAdjudicator(ctx, binding.NewConcludeExecuteMsg, req)
}

func (a *Adjudicator) withdraw(ctx context.Context, req channel.AdjudicatorReq) error {
	w := binding.NewWithdrawal(req.Params.ID(), req.Acc.Address(), a.Account())
	b := w.Bytes()
	sig, err := req.Acc.SignData(b)
	if err != nil {
		return fmt.Errorf("signing: %w", err)
	}

	msg, err := binding.NewWithdrawMsgExecute(req.Params.ID(), req.Acc.Address(), a.Account(), sig)
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}

	_, err = a.Execute(ctx, msg, nil)
	return err
}

// Progress progresses the state of a previously registered channel on-chain.
// The signatures for the old state can be nil as the state is already
// registered on the adjudicator.
func (a *Adjudicator) Progress(ctx context.Context, req channel.ProgressReq) error {
	return fmt.Errorf("unsupported")
}

// Subscribe returns an AdjudicatorEvent subscription.
//
// The context should only be used to establish the subscription. The
// framework will call Close on the subscription once the respective channel
// controller shuts down.
func (a *Adjudicator) Subscribe(ctx context.Context, ch channel.ID) (channel.AdjudicatorSubscription, error) {
	sub := NewEventSubscription(a, ch)
	return sub, nil
}
