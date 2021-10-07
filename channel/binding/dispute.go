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
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/perun-network/perun-cosmwasm-backend/channel/contract"
	"github.com/perun-network/perun-cosmwasm-backend/pkg/safecast"
	"github.com/xeipuuv/gojsonschema"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

type DisputeExecuteMsg struct {
	Dispute SignedState `json:"dispute"`
}

func NewDisputeExecuteMsg(p channel.Params, s channel.State, sigs []wallet.Sig) ([]byte, error) {
	msg := DisputeExecuteMsg{
		Dispute: makeSignedState(p, s, sigs),
	}
	return json.Marshal(msg)
}

type SignedState struct {
	Params Params `json:"params"`
	State  State  `json:"state"`
	Sigs   []Sig  `json:"sigs"`
}

func makeSignedState(p channel.Params, s channel.State, sigs []wallet.Sig) SignedState {
	params := makeParams(&p)
	state := makeState(&s)
	_sigs := makeSigs(sigs)
	return SignedState{params, state, _sigs}
}

type DisputeQueryMsg struct {
	Dispute ChannelID `json:"dispute"`
}

func NewDisputeQueryMsg(cID channel.ID) ([]byte, error) {
	msg := DisputeQueryMsg{
		Dispute: cID[:],
	}
	return json.Marshal(msg)
}

type (
	DisputeQueryResponseMsg struct {
		Concluded bool      `json:"concluded"`
		State     State     `json:"state"`
		Timeout   Timestamp `json:"timeout"`
	}

	DisputeQueryResponse struct {
		DisputeQueryResponseMsg
	}

	Timestamp = Uint64
)

var disputeReponseSchema = func() gojsonschema.JSONLoader {
	schema := contract.DisputeResponseSchema
	return gojsonschema.NewStringLoader(schema)
}()

// DecodeDisputeQueryResponse decodes a dispute query response from the given byte slice.
func DecodeDisputeQueryResponse(b []byte) (DisputeQueryResponse, error) {
	err := validateMsg(b, disputeReponseSchema)
	if err != nil {
		return DisputeQueryResponse{}, fmt.Errorf("validating: %w", err)
	}

	var resp DisputeQueryResponseMsg
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return DisputeQueryResponse{}, fmt.Errorf("unmarshalling: %w", err)
	}
	return DisputeQueryResponse{resp}, nil
}

// Timeout returns the dispute timeout.
func (r *DisputeQueryResponse) Timeout() time.Time {
	t := r.DisputeQueryResponseMsg.Timeout.Val()
	return time.Unix(0, safecast.Int64FromUint64(t))
}

// Equal returns whether receiver and argument are equal.
func (d DisputeQueryResponse) Equal(_d DisputeQueryResponse) bool {
	b, err := encodeCanonical(d)
	if err != nil {
		panic(err)
	}
	_b, err := encodeCanonical(_d)
	if err != nil {
		panic(err)
	}
	return bytes.Equal(b, _b)
}

func makeSigs(sigs []wallet.Sig) []Sig {
	_sigs := make([]Sig, len(sigs))
	for i, sig := range sigs {
		_sigs[i] = sig
	}
	return _sigs
}
