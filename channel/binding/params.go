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
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

type (
	Params struct {
		DisputeDuration Uint64        `json:"dispute_duration"`
		Nonce           ByteArray     `json:"nonce"`
		Parts           []OffIdentity `json:"participants"`
	}

	OffIdentity = ByteArray
)

func makeParams(p *channel.Params) Params {
	return Params{
		Nonce:           p.Nonce.Bytes(),
		Parts:           makeParticipants(p.Parts),
		DisputeDuration: makeUint64(p.ChallengeDuration),
	}
}

func makeParticipants(parts []wallet.Address) []OffIdentity {
	_parts := make([]OffIdentity, len(parts))
	for i, p := range parts {
		_parts[i] = p.Bytes()
	}
	return _parts
}

// Bytes returns a canonical byte representation of the object.
func (p *Params) Bytes() []byte {
	b, err := encodeCanonical(p)
	if err != nil {
		panic(err)
	}
	return b
}

func NewParams(p *channel.Params) *Params {
	_p := makeParams(p)
	return &_p
}
