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

package io

import (
	"fmt"
	"io"

	"github.com/perun-network/perun-cosmwasm-backend/pkg/safecast"
	pio "perun.network/go-perun/pkg/io"
)

// WriteBytesUint16 writes the uint16-length-encoded byte slice to the stream.
func WriteBytesUint16(w io.Writer, a []byte) error {
	l := safecast.Uint16FromInt(len(a))
	return pio.Encode(w, l, a)
}

// ReadBytesUint16 reads an uint16-length-encoded byte slice from the stream.
func ReadBytesUint16(r io.Reader, a *[]byte) error {
	var l uint16
	err := pio.Decode(r, &l)
	if err != nil {
		return fmt.Errorf("reading length: %w", err)
	}

	*a = make([]byte, l)
	return pio.Decode(r, a)
}
