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

package contract

import _ "embed"

//go:embed perun_cosmwasm.wasm
var Code []byte

//go:embed schema/init_msg.json
var InitMsgSchema string

//go:embed schema/execute_msg.json
var ExecuteMsgSchema string

//go:embed schema/query_msg.json
var QueryMsgSchema string

//go:embed schema/deposit_response.json
var DepositResponseSchema string

//go:embed schema/dispute_response.json
var DisputeResponseSchema string
