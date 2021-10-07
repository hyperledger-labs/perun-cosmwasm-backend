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

package cosmwasm

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type (
	// ContractTemplate represents a contract that can be deployed on a CosmWasm ledger.
	ContractTemplate interface {
		Code() []byte
		ValidateInitMsg([]byte) error
		ValidateExecuteMsg([]byte) error
		ValidateQueryMsg([]byte) error
	}

	// StoredContract represents a contract that is stored on the ledger.
	StoredContract interface {
		ContractTemplate
		ID() uint64
	}

	// ContractInstance represents a deployed contract instance that is ready for interaction.
	ContractInstance interface {
		StoredContract
		Address() string
	}
)

type contractTemplate struct {
	code        []byte
	initSchema  gojsonschema.JSONLoader
	execSchema  gojsonschema.JSONLoader
	querySchema gojsonschema.JSONLoader
}

func NewContractTemplate(code []byte, initSchema, execSchema, querySchema string) ContractTemplate {
	_code := code
	_initSchema := gojsonschema.NewStringLoader(initSchema)
	_execSchema := gojsonschema.NewStringLoader(execSchema)
	_querySchema := gojsonschema.NewStringLoader(querySchema)

	return &contractTemplate{
		code:        _code,
		initSchema:  _initSchema,
		execSchema:  _execSchema,
		querySchema: _querySchema,
	}
}

func (c *contractTemplate) Code() []byte {
	return c.code
}

func (c *contractTemplate) ValidateInitMsg(msg []byte) error {
	return c.validateMsg(msg, c.initSchema)
}

func (c *contractTemplate) ValidateExecuteMsg(msg []byte) error {
	return c.validateMsg(msg, c.execSchema)
}

func (c *contractTemplate) ValidateQueryMsg(msg []byte) error {
	return c.validateMsg(msg, c.querySchema)
}

// validateMsg validates the given message against the given schema.
func (*contractTemplate) validateMsg(msg []byte, schema gojsonschema.JSONLoader) error {
	_msg := gojsonschema.NewBytesLoader(msg)
	result, err := gojsonschema.Validate(schema, _msg)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return fmt.Errorf("invalid message: %v, %v", msg, result.Errors())
	}
	return nil
}

type storedContract struct {
	ContractTemplate
	id uint64
}

func NewStoredContract(c ContractTemplate, id uint64) StoredContract {
	return &storedContract{
		ContractTemplate: c,
		id:               id,
	}
}

func (c *storedContract) ID() uint64 {
	return c.id
}

type contractInstance struct {
	StoredContract
	address string
}

func NewContractInstance(c StoredContract, addr string) ContractInstance {
	return &contractInstance{
		StoredContract: c,
		address:        addr,
	}
}

func (c *contractInstance) Address() string {
	return c.address
}
