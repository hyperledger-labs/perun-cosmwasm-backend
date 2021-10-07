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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/gowebpki/jcs"
	"github.com/xeipuuv/gojsonschema"
)

func NewBigIntFromUint128Bytes(b []byte) *big.Int {
	return (&big.Int{}).SetBytes(b)
}

func NewIntFromUint128Bytes(b []byte) types.Int {
	return types.NewIntFromBigInt(NewBigIntFromUint128Bytes(b))
}

var uint128Max = new(big.Int).Lsh(big.NewInt(1), 128)

// Uint128 represents a JSON-encodable 128-bit unsigned integer.
type Uint128 struct {
	val *big.Int
}

func makeUint128(i *big.Int) Uint128 {
	if i.Sign() < 0 || i.Cmp(uint128Max) >= 0 {
		panic(fmt.Sprintf("must be in [0, 2^128]: %v", i))
	}
	_i := new(big.Int).Set(i)
	return Uint128{
		val: _i,
	}
}

func (i Uint128) Int() types.Int {
	return types.NewIntFromBigInt(i.val)
}

func (i Uint128) MarshalJSON() ([]byte, error) {
	_i := "\"" + i.val.String() + "\"" // Add quotes.
	return []byte(_i), nil
}

func (i *Uint128) UnmarshalJSON(data []byte) error {
	_data := data[1 : len(data)-1] // Strip quotes.
	_i, ok := new(big.Int).SetString(string(_data), 10)
	if !ok {
		return fmt.Errorf("failed to parse: %v", data)
	}
	*i = makeUint128(_i)
	return nil
}

type Sig = ByteArray

func validateMsg(msg []byte, schema gojsonschema.JSONLoader) error {
	_msg := gojsonschema.NewBytesLoader(msg)
	result, err := gojsonschema.Validate(schema, _msg)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return fmt.Errorf("invalid message: %v, %v", msg, result.Errors())
	}
	return nil
}

// Uint64 represents a JSON-encodable 64-bit unsigned integer.
type Uint64 uint64

func makeUint64(i uint64) Uint64 {
	return Uint64(i)
}

func (i Uint64) Val() uint64 {
	return uint64(i)
}

func (i Uint64) MarshalJSON() ([]byte, error) {
	_i := "\"" + strconv.FormatUint(uint64(i), 10) + "\"" // Add quotes.
	return []byte(_i), nil
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	_data := data[1 : len(data)-1] // Strip quotes.
	_i, err := strconv.ParseUint(string(_data), 10, 64)
	if err != nil {
		return err
	}
	*i = makeUint64(_i)
	return nil
}

// encodeCanonical creates a canonical encoding of an object.
func encodeCanonical(o interface{}) ([]byte, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return jcs.Transform(b)
}

// ByteArray represents a byte array that is JSON-encodable.
type ByteArray []byte

func (a ByteArray) MarshalJSON() ([]byte, error) {
	_a := "\"" + base64.StdEncoding.EncodeToString(a) + "\"" // Add quotes.
	return []byte(_a), nil
}

func (a *ByteArray) UnmarshalJSON(data []byte) error {
	_data := data[1 : len(data)-1] // Strip quotes.
	_a, err := base64.StdEncoding.DecodeString(string(_data))
	if err != nil {
		return err
	}
	*a = ByteArray(_a)
	return nil
}

type Hash = ByteArray
