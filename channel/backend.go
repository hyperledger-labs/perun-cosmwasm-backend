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

package channel

import (
	"crypto/sha256"
	"io"

	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

// Backend provides a set of utility functions.
type Backend struct{}

func NewBackend() *Backend {
	return &Backend{}
}

// CalcID infers the id of a channel from its parameters.
func (*Backend) CalcID(p *channel.Params) channel.ID {
	_p := binding.NewParams(p)
	b := _p.Bytes()
	return sha256.Sum256(b)
}

// Sign signs a channel's State with the given Account.
func (*Backend) Sign(a wallet.Account, s *channel.State) (wallet.Sig, error) {
	b := binding.NewState(s).Bytes()
	return a.SignData(b)
}

// Verify verifies that the provided signature on the state belongs to the
// provided address.
func (b *Backend) Verify(addr wallet.Address, s *channel.State, sig wallet.Sig) (bool, error) {
	msg := binding.NewState(s).Bytes()
	return wallet.VerifySignature(msg, sig, addr)
}

// DecodeAsset decodes an asset from a stream.
func (*Backend) DecodeAsset(r io.Reader) (channel.Asset, error) {
	var a binding.Asset
	err := a.Decode(r)
	if err != nil {
		return nil, err
	}
	return a, nil
}
