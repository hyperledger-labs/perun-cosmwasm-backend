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
	"testing"

	pkgtest "perun.network/go-perun/pkg/test"
	"perun.network/go-perun/wallet/test"
)

type NewSetupFunc func(rng *rand.Rand) *test.Setup

// TestWallet tests a wallet implementation.
func TestWallet(t *testing.T, newWalletSetup NewSetupFunc) {
	t.Run("Generic Address Test", func(t *testing.T) {
		t.Parallel()
		rng := pkgtest.Prng(t, "address")
		test.TestAddress(t, newWalletSetup(rng))
	})
	t.Run("Test Account, Wallet, and Backend", func(t *testing.T) {
		t.Parallel()
		rng := pkgtest.Prng(t, "signature")
		test.TestAccountWithWalletAndBackend(t, newWalletSetup(rng))
	})
	t.Run("Generic Signature Size Test", func(t *testing.T) {
		t.Parallel()
		rng := pkgtest.Prng(t, "signature size")
		test.GenericSignatureSizeTest(t, newWalletSetup(rng))
	})
}
