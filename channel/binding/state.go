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

package binding

import (
	perun "github.com/perun-network/perun-cosmwasm-backend/pkg/perun/channel"
	"github.com/perun-network/perun-cosmwasm-backend/pkg/safecast"
	"perun.network/go-perun/channel"
)

type State struct {
	Balances  Balances  `json:"balances"`
	ChannelID ChannelID `json:"channel_id"`
	Finalized bool      `json:"finalized"`
	Version   Uint64    `json:"version"`
}

type ChannelID = Hash

func makeState(s *channel.State) State {
	return State{
		ChannelID: s.ID[:],
		Version:   makeUint64(s.Version),
		Balances:  makeBalances(s.Allocation),
		Finalized: s.IsFinal,
	}
}

func makeBalances(a channel.Allocation) []Balance {
	if len(a.Locked) > 0 {
		panic("locked")
	}

	assets := a.Assets
	bals := a.Balances
	if len(bals) == 0 {
		panic("empty balances")
	}
	numParts := len(bals[0])
	_bals := make([]Balance, numParts)

	for i := range _bals {
		partIdx := safecast.Uint16FromInt(i)
		partBals := perun.Balances(bals).ForPart(partIdx)
		_bals[i] = makeBalance(assets, partBals)
	}
	return _bals
}

func makeBalance(assets []channel.Asset, bals []channel.Bal) Balance {
	if len(assets) != len(bals) {
		panic("length mismatch")
	}

	_bals := make(Balance, len(bals))
	for i, bal := range bals {
		denom := makeDenom(assets[i])
		amount := makeUint128(bal)
		_bals[i] = makeCoin(denom, amount)
	}
	return _bals
}

func (s *State) PerunState() *channel.State {
	var cID channel.ID
	copy(cID[:], s.ChannelID)
	return &channel.State{
		ID:         cID,
		Version:    s.Version.Val(),
		App:        channel.NoApp(),
		Allocation: Balances(s.Balances).perunAllocation(),
		Data:       channel.NoData(),
		IsFinal:    s.Finalized,
	}
}

type (
	Balances []Balance
	Balance  []Coin
)

func (b Balances) perunAllocation() channel.Allocation {
	numParts := len(b)
	numAssets := len(b[0])
	assets := make([]channel.Asset, numAssets)
	balances := make(channel.Balances, numAssets)

	for a := 0; a < numAssets; a++ {
		assets[a] = Asset(b[0][a].Denom)
		balances[a] = make([]channel.Bal, numParts)
		for p := 0; p < numParts; p++ {
			balances[a][p] = b[p][a].Amount.Int().BigInt()
		}
	}

	return channel.Allocation{
		Assets:   assets,
		Balances: balances,
		Locked:   nil,
	}
}

// Bytes returns a canonical byte representation of the object.
func (s *State) Bytes() []byte {
	b, err := encodeCanonical(s)
	if err != nil {
		panic(err)
	}
	return b
}

func NewState(s *channel.State) *State {
	_s := makeState(s)
	return &_s
}
