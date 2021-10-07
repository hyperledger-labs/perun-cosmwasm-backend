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

package test

import (
	"math/big"
	"math/rand"

	btest "github.com/perun-network/perun-cosmwasm-backend/channel/binding/test"
	"perun.network/go-perun/channel"
	ctest "perun.network/go-perun/channel/test"
)

// RandomGenerator generates randomized objects for testing.
type RandomGenerator struct {
	maxNumParts          int
	maxNumAssets         int
	maxFundingAmount     *big.Int
	maxChallengeDuration int64
}

func NewRandomGenerator(maxNumParts int, maxNumAssets int, maxFundingAmount *big.Int, maxChallengeDuration int64) *RandomGenerator {
	return &RandomGenerator{
		maxNumParts:          maxNumParts,
		maxNumAssets:         maxNumAssets,
		maxFundingAmount:     maxFundingAmount,
		maxChallengeDuration: maxChallengeDuration,
	}
}

// NewParamsAndState generates random parameters and state such that the participants' accounts are in the wallet.
func (r *RandomGenerator) NewParamsAndState(rng *rand.Rand, opts ...ctest.RandomOpt) (*channel.Params, *channel.State) {
	numParts := 2 + rng.Intn(r.maxNumParts-2)
	numAssets := 1 + rng.Intn(r.maxNumAssets-1)
	assets := NewRandomAssets(rng, numAssets)
	challengeDuration := uint64(rng.Int63n(r.maxChallengeDuration))

	_opts := make([]ctest.RandomOpt, 0)
	_opts = append(
		_opts,
		ctest.WithNumParts(numParts),
		ctest.WithNumLocked(0),
		ctest.WithBalancesInRange(big.NewInt(0), r.maxFundingAmount),
		ctest.WithAssets(assets...),
		ctest.WithChallengeDuration(challengeDuration),
	)
	_opts = append(_opts, opts...)

	return ctest.NewRandomParamsAndState(rng, _opts...)
}

// NewRandomAssets generates a set of randpm assets.
func NewRandomAssets(rng *rand.Rand, n int) []channel.Asset {
	a := make([]channel.Asset, n)
	for i := range a {
		a[i] = btest.NewRandomAsset(rng)
	}
	return a
}
