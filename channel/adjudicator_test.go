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

package channel_test

import (
	"context"
	"math/big"
	"math/rand"
	"testing"

	bchannel "github.com/perun-network/perun-cosmwasm-backend/channel"
	"github.com/perun-network/perun-cosmwasm-backend/channel/test"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	"github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm/simulation"
	ptest "github.com/perun-network/perun-cosmwasm-backend/pkg/perun/channel/test"
	"perun.network/go-perun/channel"
	ctest "perun.network/go-perun/channel/test"
	pkgtest "perun.network/go-perun/pkg/test"
	"perun.network/go-perun/wallet"
	wtest "perun.network/go-perun/wallet/test"
)

// TestAdjudicator tests the adjudicator.
func TestAdjudicator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	rng := pkgtest.Prng(t)
	c, contract := test.NewTestClientWithContract(ctx, t)
	c.StartTicking(blockTick, simChainTick)
	defer c.StopTicking()
	a := newAdjudicatorSetup(c, contract)
	ptest.TestAdjudicatorWithSubscription(ctx, t, rng, a)
}

type adjudicatorSetup struct {
	c        *simulation.Client
	contract client.ContractInstance
	adj      channel.Adjudicator
	r        *test.RandomGenerator
	w        wtest.Wallet
}

func newAdjudicatorSetup(c *simulation.Client, contract client.ContractInstance) *adjudicatorSetup {
	opt := bchannel.AdjudicatorPollingIntervalOpt(polling)
	adj := bchannel.NewAdjudicator(c, contract, c.Account(), opt)
	return &adjudicatorSetup{
		c:        c,
		contract: contract,
		adj:      &testAdjudicator{adj},
		r:        test.NewRandomGenerator(maxNumParts, maxNumAssets, big.NewInt(maxFundingAmount), int64(maxChallengeDuration.Seconds())),
		w:        wtest.RandomWallet(),
	}
}

type testAdjudicator struct {
	*bchannel.Adjudicator
}

// Progress is a no-op because app channels are not supported yet.
func (a *testAdjudicator) Progress(ctx context.Context, req channel.ProgressReq) error {
	return nil
}

func (a *adjudicatorSetup) Adjudicator() channel.Adjudicator {
	return a.adj
}

func (a *adjudicatorSetup) NewFundedChannel(ctx context.Context, rng *rand.Rand) (channel.Params, channel.State) {
	opts := []ctest.RandomOpt{ctest.WithoutApp(), ctest.WithIsFinal(false), ctest.WithVersion(0)}
	params, state := a.r.NewParamsAndState(rng, opts...)
	a.fund(ctx, params, state)
	return *params, *state
}

// SignState signs the channel state on behalf all participants.
func (a *adjudicatorSetup) SignState(state *channel.State, parts []wallet.Address) []wallet.Sig {
	sigs := make([]wallet.Sig, len(parts))
	for i, p := range parts {
		acc, err := a.w.Unlock(p)
		if err != nil {
			panic(err)
		}
		sigs[i], err = channel.Sign(acc, state)
		if err != nil {
			panic(err)
		}
	}
	return sigs
}

// Account returns the account for an address.
func (a *adjudicatorSetup) Account(addr wallet.Address) wallet.Account {
	acc, err := a.w.Unlock(addr)
	if err != nil {
		panic(err)
	}
	return acc
}

// fund funds the specified channel.
func (a *adjudicatorSetup) fund(ctx context.Context, params *channel.Params, state *channel.State) {
	requests := make([]*channel.FundingReq, len(params.Parts))
	for i := range requests {
		requests[i] = newFundingRequest(ctx, params, state, channel.Index(i), a.c)
	}

	opt := bchannel.FunderPollingIntervalOpt(polling)
	f := bchannel.NewFunder(a.c, a.contract, a.c.Account(), opt)
	err := fundAll(ctx, f, requests)
	if err != nil {
		panic(err)
	}
}

// fundAll executes a set of funding requests.
func fundAll(ctx context.Context, f channel.Funder, requests []*channel.FundingReq) error {
	errs := make(chan error, len(requests))
	for _, req := range requests {
		go func(req channel.FundingReq) { errs <- f.Fund(ctx, req) }(*req)
	}
	for range requests {
		select {
		case err := <-errs:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
