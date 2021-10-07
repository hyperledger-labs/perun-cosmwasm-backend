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
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	perun "github.com/perun-network/perun-cosmwasm-backend/pkg/perun/channel"
	"perun.network/go-perun/channel"
)

// Funder provides methods for funding a channel.
type Funder struct {
	*contractClient
	polling time.Duration
}

type FunderOpt func(*Funder)

func FunderPollingIntervalOpt(d time.Duration) FunderOpt {
	return func(f *Funder) {
		f.polling = d
	}
}

func NewFunder(c client.Client, contract client.ContractInstance, acc types.AccAddress, opts ...FunderOpt) *Funder {
	f := &Funder{
		contractClient: newContractClient(c, contract, acc),
		polling:        defaultPollingInterval,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Fund deposits funds according to the specified funding request and waits until the funding is complete.
func (f *Funder) Fund(ctx context.Context, req channel.FundingReq) error {
	_req := (*fundingReq)(&req)
	fID, err := _req.ID()
	if err != nil {
		return fmt.Errorf("creating funding ID: %w", err)
	}

	funds := _req.Funds()
	err = f.deposit(ctx, fID, funds)
	if err != nil {
		return fmt.Errorf("depositing: %w", err)
	}
	return f.awaitFundingComplete(ctx, _req)
}

type fundingReq channel.FundingReq

func (r *fundingReq) ID() (binding.FundingID, error) {
	return r.IDForPart(r.Idx)
}

func (r *fundingReq) IDForPart(i channel.Index) (binding.FundingID, error) {
	return binding.CalcFundingID(r.Params.ID(), r.Params.Parts[i])
}

func (r *fundingReq) Funds() types.Coins {
	bals := perun.Balances(r.Agreement).ForPart(r.Idx)
	return binding.MakeCoins(r.State.Assets, bals)
}

func (r *fundingReq) TotalFunds() types.Coins {
	return binding.MakeCoins(r.State.Assets, r.State.Allocation.Sum())
}

// deposit sends a funding transaction.
func (f *Funder) deposit(ctx context.Context, fID binding.FundingID, funds types.Coins) error {
	msg, err := binding.NewDepositExecuteMsg(fID)
	if err != nil {
		return err
	}

	_, err = f.Execute(ctx, msg, funds)
	return err
}

// awaitFundingComplete blocks until the funding of the specified channel is complete.
func (f *Funder) awaitFundingComplete(ctx context.Context, req *fundingReq) error {
	total := req.TotalFunds()
	for {
		funded := types.NewCoins()
		for i := range req.Params.Parts {
			_funded, err := f.queryDeposit(ctx, req, channel.Index(i))
			if err != nil {
				log.Printf("Warning: Error querying deposit: %v\n", err)
			}
			funded = funded.Add(_funded...)
		}

		if funded.IsAllGTE(total) {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(f.polling):
		}
	}

}

// queryDeposit queries the current deposit state for the given channel participant.
func (f *Funder) queryDeposit(ctx context.Context, req *fundingReq, part channel.Index) (binding.DepositQueryResponse, error) {
	fID, err := req.IDForPart(part)
	if err != nil {
		return nil, err
	}
	msg, err := binding.NewDepositQueryMsg(fID)
	if err != nil {
		return nil, err
	}

	resp, err := f.Query(ctx, msg)
	if err != nil {
		return nil, err
	}

	return binding.DecodeDepositQueryResponse(resp.Data)
}
