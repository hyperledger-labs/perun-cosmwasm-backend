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
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

type WithdrawExecuteMsg struct {
	Withdraw SignedWithdrawal `json:"withdraw"`
}

func NewWithdrawMsgExecute(p channel.ID, part wallet.Address, receiver types.Address, sig wallet.Sig) ([]byte, error) {
	w := makeWithdrawal(p, part, receiver)
	msg := WithdrawExecuteMsg{
		Withdraw: makeSignedWithdrawal(w, sig),
	}
	return json.Marshal(msg)
}

type SignedWithdrawal struct {
	Withdrawal Withdrawal `json:"withdrawal"`
	Sig        Sig        `json:"sig"`
}

func makeSignedWithdrawal(w Withdrawal, sig wallet.Sig) SignedWithdrawal {
	return SignedWithdrawal{w, sig}
}

type Withdrawal struct {
	ChannelId ByteArray   `json:"channel_id"`
	Part      OffIdentity `json:"part"`
	Receiver  string      `json:"receiver"`
}

func makeWithdrawal(p channel.ID, part wallet.Address, receiver types.Address) Withdrawal {
	return Withdrawal{
		ChannelId: p[:],
		Part:      part.Bytes(),
		Receiver:  receiver.String(),
	}
}

// Bytes returns a canonical byte representation of the object.
func (w *Withdrawal) Bytes() []byte {
	b, err := json.Marshal(w)
	if err != nil {
		panic(err)
	}
	return b
}

func NewWithdrawal(p channel.ID, part wallet.Address, receiver types.Address) *Withdrawal {
	w := makeWithdrawal(p, part, receiver)
	return &w
}
