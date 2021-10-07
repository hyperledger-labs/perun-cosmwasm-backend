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
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/types"
	bio "github.com/perun-network/perun-cosmwasm-backend/pkg/io"
	"perun.network/go-perun/channel"
)

// MakeCoins translates assets and balances into coins.
func MakeCoins(assets []channel.Asset, bals []channel.Bal) types.Coins {
	if len(assets) != len(bals) {
		panic("mismatch")
	}
	coinsList := make([]types.Coin, len(assets))
	for i, bal := range bals {
		denom := makeDenom(assets[i])
		amount := types.NewIntFromBigInt(bal)
		coinsList[i] = types.NewCoin(denom, amount)
	}
	return types.NewCoins(coinsList...)
}

// Denom represents an asset denomination.
type Denom = string

func makeDenom(a channel.Asset) Denom {
	denom, ok := a.(Asset)
	if !ok {
		panic(fmt.Sprintf("invalid type: %T", a))
	}
	return Denom(denom)
}

// Asset represents an asset.
type Asset string

// Decode decodes an asset from a reader.
func (a *Asset) Decode(r io.Reader) error {
	var b []byte
	err := bio.ReadBytesUint16(r, &b)
	if err != nil {
		return err
	}
	*a = Asset(b)
	return nil
}

// Encodes an asset onto a writer.
func (a Asset) Encode(w io.Writer) error {
	return bio.WriteBytesUint16(w, []byte(a))
}
