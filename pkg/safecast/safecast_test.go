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

package safecast_test

import (
	"math/rand"
	"testing"

	"github.com/perun-network/perun-cosmwasm-backend/pkg/safecast"
)

func TestUint16FromInt(t *testing.T) {
	// Test valid input: Correct result, no panic.
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("casting valid value should not panic: %v", err)
			}
		}()

		i := rand.Intn(1 << 16)
		_i := safecast.Uint16FromInt(i)
		if i != int(_i) {
			t.Errorf("values should be equal: %v, %v", i, _i)
		}
	}()

	// Test invalid input: Panic.
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("casting invalid value should panic")
			}
		}()

		i := 1 << 16
		_ = safecast.Uint16FromInt(i)
	}()
}

func TestInt64FromUint64(t *testing.T) {
	// Test valid input: Correct result, no panic.
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("casting valid value should not panic: %v", err)
			}
		}()

		i := uint64(rand.Int63())
		_i := safecast.Int64FromUint64(i)
		if i != uint64(_i) {
			t.Errorf("values should be equal: %v, %v", i, _i)
		}
	}()

	// Test invalid input: Panic.
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("casting invalid value should panic")
			}
		}()

		i := uint64(1<<63) + 1
		_ = safecast.Int64FromUint64(i)
	}()
}
