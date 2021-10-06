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
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	bwallet "github.com/perun-network/perun-cosmwasm-backend/wallet"
	"perun.network/go-perun/wallet"
)

// Wallet provides methods for wallet testing.
type Wallet struct {
	*bwallet.Wallet
	counter int
}

// NewWallet creates a new test wallet.
func NewWallet() *Wallet {
	kr := newKeyring()
	return &Wallet{
		Wallet:  bwallet.NewWallet(kr),
		counter: 0,
	}
}

func newKeyring() keyring.Keyring {
	kr, err := keyring.New("", keyring.BackendMemory, "", nil)
	if err != nil {
		panic(err)
	}
	return kr
}

// NewRandomAccount generates a new account.
func (w *Wallet) NewRandomAccount(rng *rand.Rand) wallet.Account {
	pwd := ""
	w.counter++
	name := fmt.Sprintf("Account%d", w.counter)
	acc, err := w.Wallet.NewAccount(rng, name, pwd)
	if err != nil {
		panic(err)
	}
	return acc
}
