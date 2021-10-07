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

package channel_test

import (
	"testing"

	"github.com/perun-network/perun-cosmwasm-backend/pkg/perun/channel/test"
	ctest "perun.network/go-perun/channel/test"
)

// TestBackend tests the backend.
func TestBackend(t *testing.T) {
	opts := []ctest.RandomOpt{
		ctest.WithNumLocked(0), // locked funds are not supported yet
		ctest.WithoutApp(),     // app and data are not supported yet
	}
	test.TestChannelBackend(t, opts...)
}
