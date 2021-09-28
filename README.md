<h1 align="center"><br>
    <a href="https://perun.network/"><img src=".assets/perun.png" alt="Perun" width="30%"></a>
<br></h1>

<h2 align="center">Perun CosmWasm: Backend</h2>

<p align="center">
  <a href="https://www.apache.org/licenses/LICENSE-2.0.txt"><img src="https://img.shields.io/badge/license-Apache%202-blue" alt="License: Apache 2.0"></a>
  </a>
  <a href="https://github.com/perun-network/perun-cosmwasm-backend/actions/workflows/rust.yml"><img src="https://github.com/perun-network/perun-cosmwasm-backend/actions/workflows/go.yml/badge.svg?branch=main" alt="CI status"></a>
  </a>
</p>

This repository contains an implementation of a [CosmWasm](https://www.cosmwasm.com) backend for the [go-perun](https://perun.network/) state channel library.

## Organization
* `channel/`: Implementation of the `go-perun/channel` interfaces.
  * `contract/`: Contains the compiled contract and JSON Schema files from [perun-cosmwasm-contract](https://github.com/perun-network/perun-cosmwasm-contract).
* `client/`: End-to-end tests.
* `wallet/`: Implementation of the `go-perun/wallet` interfaces.

## Dependencies

[Golang](https://golang.org) version >= 1.17 must be installed.

## Usage

The end-to-end tests in package `client` demonstrate how to setup and run a Perun client on a CosmWasm network.

To run all unit and end-to-end tests, open a Terminal at the repository root directory and run

```sh
go test ./...
```

## Copyright

Copyright 2021 PolyCrypt GmbH.

Use of the source code is governed by the Apache 2.0 license that can be found in the [LICENSE file](LICENSE).
