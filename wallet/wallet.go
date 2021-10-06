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

package wallet

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
	"perun.network/go-perun/wallet"
)

const mnemonicEntropySize = 128 / 8

// Wallet is an implementation of wallet.Wallet based on keyring.
type Wallet struct {
	kr keyring.Keyring
	mu sync.Mutex
}

// NewWallet creates a new wallet based on the specified keyring.
func NewWallet(kr keyring.Keyring) *Wallet {
	return &Wallet{
		kr: kr,
		mu: sync.Mutex{},
	}
}

// NewAccount creates a new account and returns it.
func (w *Wallet) NewAccount(rng *rand.Rand, name string, pwd string) (wallet.Account, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	kr := w.kr
	algos, _ := kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return nil, err
	}

	m, err := w.newMnemonic(rng)
	if err != nil {
		return nil, err
	}

	hdPath := ""
	acc, err := kr.NewAccount(name, m, pwd, hdPath, algo)
	if err != nil {
		return nil, err
	}
	return &Account{
		kr:   kr,
		addr: NewAddress(acc.GetPubKey()),
	}, nil
}

func (w *Wallet) newMnemonic(rand *rand.Rand) (string, error) {
	var entropy [mnemonicEntropySize]byte
	_, err := rand.Read(entropy[:])
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy[:])
}

// Unlock unlocks an account from the wallet.
func (w *Wallet) Unlock(addr wallet.Address) (wallet.Account, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.contains(addr) {
		return nil, fmt.Errorf("unknown account: %v", addr)
	}

	a, err := AsAddr(addr)
	if err != nil {
		return nil, err
	}

	return &Account{
		kr:   w.kr,
		addr: a,
	}, nil
}

// LockAll implements wallet.LockAll. It is noop.
func (w *Wallet) LockAll() {}

// IncrementUsage implements wallet.Wallet. It is a noop.
func (w *Wallet) IncrementUsage(a wallet.Address) {}

// DecrementUsage implements wallet.Wallet. It is a noop.
func (w *Wallet) DecrementUsage(a wallet.Address) {}

func (w *Wallet) contains(addr wallet.Address) bool {
	a, err := AsAddr(addr)
	if err != nil {
		return false
	}

	_, err = w.kr.KeyByAddress(a.CosmAddr())
	return err == nil
}
