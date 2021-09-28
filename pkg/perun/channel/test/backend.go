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
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"perun.network/go-perun/channel"
	channeltest "perun.network/go-perun/channel/test"
	pkgtest "perun.network/go-perun/pkg/test"
	"perun.network/go-perun/wallet"
	wallettest "perun.network/go-perun/wallet/test"
)

// TestChannelBackend tests the global channel backend.
func TestChannelBackend(t *testing.T, opts ...channeltest.RandomOpt) {
	setup := newChannelSetup(pkgtest.Prng(t), opts...)
	channeltest.GenericBackendTest(t, setup)
	decodeAssetTest(t, setup)
}

func newChannelSetup(rng *rand.Rand, opts ...channeltest.RandomOpt) *channeltest.Setup {
	opts1 := append(opts, channeltest.WithIsFinal(false))
	opts2 := append(opts, channeltest.WithIsFinal(true))
	params, state := channeltest.NewRandomParamsAndState(rng, opts1...)
	params2, state2 := channeltest.NewRandomParamsAndState(rng, opts2...)

	createAddr := func() wallet.Address {
		return wallettest.NewRandomAddress(rng)
	}

	return &channeltest.Setup{
		Params:        params,
		Params2:       params2,
		State:         state,
		State2:        state2,
		Account:       wallettest.NewRandomAccount(rng),
		RandomAddress: createAddr,
	}
}

func decodeAssetTest(t *testing.T, s *channeltest.Setup) {
	assets := s.State.Assets
	_assets := make([]channel.Asset, len(assets))
	for i, a := range assets {
		var buf bytes.Buffer
		err := a.Encode(&buf)
		assert.NoError(t, err, "encoding")
		_assets[i], err = channel.DecodeAsset(&buf)
		assert.NoError(t, err, "decoding")
	}
	err := channel.AssetsAssertEqual(assets, _assets)
	assert.NoError(t, err, "equality")
}
