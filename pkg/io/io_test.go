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

package io_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/perun-network/perun-cosmwasm-backend/pkg/io"
)

func TestReadWriteBytesUint16(t *testing.T) {
	// Generate random byte slice with uint16 length.
	l := rand.Intn(1 << 16)
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = io.WriteBytesUint16(&buf, b)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}

	_b := make([]byte, len(b))
	err = io.ReadBytesUint16(&buf, &_b)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	if !bytes.Equal(b, _b) {
		t.Errorf("no equal: %v, %v", b, _b)
	}
}
