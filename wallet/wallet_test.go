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

package wallet_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	perun "github.com/perun-network/perun-cosmwasm-backend/pkg/perun/wallet/test"
	"github.com/perun-network/perun-cosmwasm-backend/wallet"
	bwallet "github.com/perun-network/perun-cosmwasm-backend/wallet"
	wallettest "github.com/perun-network/perun-cosmwasm-backend/wallet/test"
	"perun.network/go-perun/wallet/test"
)

func TestWallet(t *testing.T) {
	perun.TestWallet(t, newWalletSetup)
}

func newWalletSetup(rng *rand.Rand) *test.Setup {
	b := bwallet.NewBackend()
	w := test.Wallet(wallettest.NewWallet())
	accountA := w.NewRandomAccount(rng)
	accountB := w.NewRandomAccount(rng)
	data := make([]byte, rng.Intn(256))
	_, err := rng.Read(data)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = accountB.Address().Encode(&buf)
	if err != nil {
		panic(err)
	}
	addrEncoded := buf.Bytes()

	return &test.Setup{
		Wallet:          w,
		AddressInWallet: accountA.Address(),
		ZeroAddress:     wallet.NewAddress(&secp256k1.PubKey{Key: make([]byte, 33)}),
		AddressEncoded:  addrEncoded,
		Backend:         b,
		DataToSign:      data,
	}
}
