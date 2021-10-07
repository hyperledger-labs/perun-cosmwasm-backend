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
	"time"

	bchannel "github.com/perun-network/perun-cosmwasm-backend/channel"
	bchanneltest "github.com/perun-network/perun-cosmwasm-backend/channel/binding/test"
	bwallet "github.com/perun-network/perun-cosmwasm-backend/wallet"
	"github.com/perun-network/perun-cosmwasm-backend/wallet/test"
	"perun.network/go-perun/channel"
	channeltest "perun.network/go-perun/channel/test"
	"perun.network/go-perun/wallet"
	wallettest "perun.network/go-perun/wallet/test"
)

const polling = 500 * time.Millisecond          // The interval at which a subscription polls for new state.
const maxNumParts = 16                          // The maxiumum number of participants in a channel.
const maxNumAssets = 8                          // The maximum number of assets used in a channel.
const maxFundingAmount = 2048                   // The maximum amount of channel funding required per participant.
const maxChallengeDuration = 3600 * time.Second // The maximum channel challenge duration.
const testTimeout = 30 * time.Second            // The duration after a test times out.
const blockTick = 10 * time.Millisecond         // The interval at which the simulated blockchain ticks.
const simChainTick = 60 * time.Second           // The amount of time that is added each blockchain tick.

// init sets the global variables for testing.
func init() {
	channel.SetBackend(bchannel.NewBackend())
	channeltest.SetRandomizer(bchanneltest.NewRandomizer())
	wallettest.SetRandomizer(test.NewRandomizer())
	wallet.SetBackend(bwallet.NewBackend())
}
