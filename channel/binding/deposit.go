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
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/perun-network/perun-cosmwasm-backend/channel/contract"
	"github.com/xeipuuv/gojsonschema"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

type FundingID = Hash

type DepositExecuteMsg struct {
	Deposit FundingID `json:"deposit"`
}

func NewDepositExecuteMsg(fID FundingID) ([]byte, error) {
	msg := DepositExecuteMsg{
		Deposit: fID,
	}
	return json.Marshal(msg)
}

func NewDepositQueryMsg(fID FundingID) ([]byte, error) {
	return NewDepositExecuteMsg(fID)
}

type (
	DepositQueryResponse = types.Coins

	DepositResponseMsg []Coin
)

var depositReponseSchema = func() gojsonschema.JSONLoader {
	schema := contract.DepositResponseSchema
	return gojsonschema.NewStringLoader(schema)
}()

// DecodeDepositQueryResponse decodes a deposit query response from the
// given byte slice.
func DecodeDepositQueryResponse(b []byte) (DepositQueryResponse, error) {
	err := validateMsg(b, depositReponseSchema)
	if err != nil {
		return nil, fmt.Errorf("validating: %w", err)
	}

	var coins DepositResponseMsg
	err = json.Unmarshal(b, &coins)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	coinsList := make([]types.Coin, len(coins))
	for i, coin := range coins {
		denom := coin.Denom
		amount := coin.Amount.Int()
		coinsList[i] = types.NewCoin(denom, amount)
	}
	return types.NewCoins(coinsList...), nil
}

type Coin struct {
	Denom  string  `json:"denom"`
	Amount Uint128 `json:"amount"`
}

func makeCoin(d string, a Uint128) Coin {
	return Coin{
		Denom:  d,
		Amount: a,
	}
}

type funding struct {
	Channel ChannelID   `json:"channel"`
	Part    OffIdentity `json:"part"`
}

// CalcFundingID calculates the funding ID for a channel participant.
func CalcFundingID(ch channel.ID, addr wallet.Address) (FundingID, error) {
	f := funding{
		Channel: ch[:],
		Part:    addr.Bytes(),
	}

	b, err := encodeCanonical(f)
	if err != nil {
		return nil, fmt.Errorf("encoding: %w", err)
	}

	h := sha256.Sum256(b)
	return h[:], nil
}
