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
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/lucasjones/reggen"
	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	"perun.network/go-perun/channel"
)

// NewRandomAsset returns a new random Asset.
func NewRandomAsset(rng *rand.Rand) binding.Asset {
	regex := types.DefaultCoinDenomRegex()
	denom, err := reggen.Generate(regex, 16)
	if err != nil {
		panic(err)
	}
	return binding.Asset(denom)
}

type Randomizer struct{}

func NewRandomizer() *Randomizer {
	return &Randomizer{}
}

func (Randomizer) NewRandomAsset(rng *rand.Rand) channel.Asset {
	return NewRandomAsset(rng)
}
