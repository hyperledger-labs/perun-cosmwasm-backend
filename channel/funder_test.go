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
	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	"github.com/perun-network/perun-cosmwasm-backend/channel/test"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	"github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm/simulation"
	pchannel "github.com/perun-network/perun-cosmwasm-backend/pkg/perun/channel"
	ptest "github.com/perun-network/perun-cosmwasm-backend/pkg/perun/channel/test"
	"perun.network/go-perun/channel"
	ctest "perun.network/go-perun/channel/test"
	pkgtest "perun.network/go-perun/pkg/test"
	wtest "perun.network/go-perun/wallet/test"
)

// TestFunder tests the funder.
func TestFunder(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	rng := pkgtest.Prng(t)
	c, contract := test.NewTestClientWithContract(ctx, t)

	numParts := 2 + rng.Intn(maxNumParts-2)
	funders := make([]channel.Funder, numParts)
	for i := range funders {
		opt := bchannel.FunderPollingIntervalOpt(polling)
		funders[i] = bchannel.NewFunder(c, contract, c.Account(), opt)
	}

	_f := &funder{
		client:   c,
		contract: contract,
		w:        wtest.RandomWallet(),
		r:        test.NewRandomGenerator(maxNumParts, maxNumAssets, big.NewInt(maxFundingAmount), int64(maxChallengeDuration.Seconds())),
		funders:  funders,
	}
	ptest.TestFunder(ctx, t, rng, _f)
}

// funder represents a funder for testing.
type funder struct {
	client   *simulation.Client
	contract client.ContractInstance
	w        wtest.Wallet
	r        *test.RandomGenerator
	funders  []channel.Funder
}

func (f *funder) Funders() []channel.Funder {
	return f.funders
}

func (f *funder) NewFundingRequests(ctx context.Context, t *testing.T, rng *rand.Rand) []channel.FundingReq {
	numParts := len(f.funders)
	params, state := f.r.NewParamsAndState(rng, ctest.WithNumParts(numParts))

	requests := make([]channel.FundingReq, numParts)
	for i := range f.funders {
		requests[i] = *newFundingRequest(ctx, params, state, channel.Index(i), f.client)
	}
	return requests
}

// newFundingRequest returns a funding request for the specified participant and ensures that the corresponding account has sufficient funds available.
func newFundingRequest(ctx context.Context, params *channel.Params, state *channel.State, idx channel.Index, client *simulation.Client) *channel.FundingReq {
	req := &channel.FundingReq{
		Params:    params,
		State:     state,
		Idx:       idx,
		Agreement: state.Balances,
	}

	coins := binding.MakeCoins(state.Assets, pchannel.Balances(state.Balances).ForPart(idx))
	err := client.AddCoins(ctx, client.Account(), coins)
	if err != nil {
		panic(err)
	}

	return req
}
