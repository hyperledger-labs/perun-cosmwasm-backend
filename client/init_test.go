package client_test

import (
	bchannel "github.com/perun-network/perun-cosmwasm-backend/channel"
	bchanneltest "github.com/perun-network/perun-cosmwasm-backend/channel/binding/test"
	bwallet "github.com/perun-network/perun-cosmwasm-backend/wallet"
	"github.com/perun-network/perun-cosmwasm-backend/wallet/test"
	"github.com/sirupsen/logrus"
	"perun.network/go-perun/channel"
	channeltest "perun.network/go-perun/channel/test"
	plogrus "perun.network/go-perun/log/logrus"
	"perun.network/go-perun/wallet"
	wallettest "perun.network/go-perun/wallet/test"
)

func init() {
	channel.SetBackend(bchannel.NewBackend())
	channeltest.SetRandomizer(bchanneltest.NewRandomizer())
	wallettest.SetRandomizer(test.NewRandomizer())
	wallet.SetBackend(bwallet.NewBackend())

	plogrus.Set(logrus.InfoLevel, &logrus.TextFormatter{ForceColors: true})
}
