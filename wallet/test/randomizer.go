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

	"perun.network/go-perun/wallet"
	wallettest "perun.network/go-perun/wallet/test"
)

// Randomizer provides methods for generating randomized test input.
type Randomizer struct {
	w *Wallet
}

func NewRandomizer() *Randomizer {
	return &Randomizer{
		w: NewWallet(),
	}
}

// NewRandomAddress returns a new random address generated from the
// passed rng.
func (r *Randomizer) NewRandomAddress(rng *rand.Rand) wallet.Address {
	return r.RandomWallet().NewRandomAccount(rng).Address()
}

// RandomWallet returns a fixed random wallet that is part of the
// randomizer's state. It will be used to generate accounts with
// NewRandomAccount.
func (r *Randomizer) RandomWallet() wallettest.Wallet {
	return r.w
}

// NewWallet returns a fresh, temporary Wallet that doesn't hold any
// accounts yet.
func (r *Randomizer) NewWallet() wallettest.Wallet {
	return NewWallet()
}
