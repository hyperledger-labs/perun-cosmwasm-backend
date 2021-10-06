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

package wallet

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"perun.network/go-perun/wallet"
)

// Account represents an account held in the HD wallet.
type Account struct {
	kr   keyring.Keyring
	addr *Address
}

// Address returns the address of this account.
func (a *Account) Address() wallet.Address {
	return a.addr
}

// SignData is used to sign data with this account.
func (a *Account) SignData(data []byte) ([]byte, error) {
	sig, _, err := a.kr.SignByAddress(a.addr.CosmAddr(), data)
	return sig, err
}
