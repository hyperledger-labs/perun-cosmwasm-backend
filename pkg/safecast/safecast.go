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

package safecast

// Uint16FromInt converts an int into an uint16 and panics if the value overflows.
func Uint16FromInt(a int) uint16 {
	b := uint16(a)
	if a != int(b) {
		panic("unsafe")
	}
	return b
}

// Int64FromUint64 converts an uint64 into an int64 and panics if the value overflows.
func Int64FromUint64(a uint64) int64 {
	b := int64(a)
	if b < 0 {
		panic("unsafe")
	}
	return b
}
