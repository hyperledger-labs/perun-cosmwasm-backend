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
	"io"

	pio "perun.network/go-perun/pkg/io"
	"perun.network/go-perun/wallet"
)

// Backend provides utility functions.
type Backend struct{}

func NewBackend() *Backend {
	return &Backend{}
}

// DecodeAddress reads and decodes an address from an io.Writer
func (*Backend) DecodeAddress(r io.Reader) (wallet.Address, error) {
	var a Address
	err := a.Decode(r)
	return &a, err
}

const SigntuareLength = 64

// DecodeSig reads a signature from the provided stream.
func (*Backend) DecodeSig(r io.Reader) (wallet.Sig, error) {
	sig := make([]byte, SigntuareLength)
	return sig, pio.Decode(r, &sig)
}

// VerifySignature verifies if this signature was signed by this address.
func (*Backend) VerifySignature(msg []byte, sig wallet.Sig, addr wallet.Address) (bool, error) {
	a, err := AsAddr(addr)
	if err != nil {
		return false, err
	}

	return a.PubKey.VerifySignature(msg, sig), nil
}
