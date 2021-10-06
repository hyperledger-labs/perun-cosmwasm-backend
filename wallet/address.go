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
	"bytes"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	cio "github.com/perun-network/perun-cosmwasm-backend/pkg/io"
	"perun.network/go-perun/wallet"
)

// Address represents an account address.
type Address struct {
	ctypes.PubKey
}

// NewAddress creates a new address.
func NewAddress(pk ctypes.PubKey) *Address {
	a := Address{pk}
	return &a
}

// CosmAddr returns the cosmos address.
func (a *Address) CosmAddr() types.Address {
	return types.AccAddress(a.Address().Bytes())
}

// AsAddr attempts to cast a generic address into a specific address.
func AsAddr(addr wallet.Address) (*Address, error) {
	a, ok := addr.(*Address)
	if !ok {
		return nil, fmt.Errorf("invalid type: %T", addr)
	}
	return a, nil
}

// Encode writes the object to a stream.
func (a *Address) Encode(w io.Writer) error {
	pk, ok := a.PubKey.(*secp256k1.PubKey)
	if !ok {
		panic(fmt.Sprintf("wrong type: %T", a.PubKey))
	}
	return cio.WriteBytesUint16(w, pk.Key)
}

// Decode reads an object from a stream.
func (a *Address) Decode(r io.Reader) error {
	var b []byte
	err := cio.ReadBytesUint16(r, &b)
	if err != nil {
		return fmt.Errorf("reading byte stream: %w", err)
	}

	pk := secp256k1.PubKey{
		Key: b,
	}
	*a = *NewAddress(&pk)
	return nil
}

// Bytes returns the representation of the address as byte slice.
func (a *Address) Bytes() []byte {
	return a.PubKey.Bytes()
}

// String converts this address to a string.
func (a *Address) String() string {
	return types.AccAddress(a.PubKey.Address().Bytes()).String()
}

// Equals returns whether the two addresses are equal. The implementation
// must be equivalent to checking `Address.Cmp(Address) == 0`.
func (a *Address) Equals(b wallet.Address) bool {
	_b, err := AsAddr(b)
	if err != nil {
		return false
	}
	return a.PubKey.Equals(_b.PubKey)
}

// Cmp compares the byte representation of two addresses. For `a.Cmp(b)`
// returns -1 if a < b, 0 if a == b, 1 if a > b.
func (a Address) Cmp(b wallet.Address) int {
	return bytes.Compare(a.Bytes(), b.Bytes())
}
