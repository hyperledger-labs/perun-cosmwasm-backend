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

package test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/perun-network/perun-cosmwasm-backend/channel/contract"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	"github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm/simulation"
	"github.com/stretchr/testify/require"
)

// NewTestClientWithContract creates a new test client and deploys the Perun contract.
func NewTestClientWithContract(ctx context.Context, t *testing.T) (*simulation.Client, client.ContractInstance) {
	c := simulation.NewTestClient(t)

	contractTemplate := client.NewContractTemplate(
		contract.Code,
		contract.InitMsgSchema,
		contract.ExecuteMsgSchema,
		contract.QueryMsgSchema,
	)

	storedContract, err := c.StoreContractTemplate(ctx, contractTemplate)
	require.NoError(t, err, "store contract")

	initMsg, err := json.Marshal(struct{}{})
	require.NoError(t, err, "create init message")

	contractInstance, _, err := c.InstantiateStoredContract(ctx, storedContract, initMsg, types.Coins{})
	require.NoError(t, err, "init contract")

	return c, contractInstance
}
