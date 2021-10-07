package client_test

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	bchannel "github.com/perun-network/perun-cosmwasm-backend/channel"
	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	chtest "github.com/perun-network/perun-cosmwasm-backend/channel/test"

	"perun.network/go-perun/channel"
	pchtest "perun.network/go-perun/channel/test"
	"perun.network/go-perun/client"
	ctest "perun.network/go-perun/client/test"
	pkgtest "perun.network/go-perun/pkg/test"
	wtest "perun.network/go-perun/wallet/test"
	"perun.network/go-perun/wire"
	wiretest "perun.network/go-perun/wire/test"
)

const defaultContextTimeout = 5 * time.Second
const roleOperationTimeout = 5 * time.Second
const challengeDuration = 60
const polling = 100 * time.Millisecond
const blockTick = 100 * time.Millisecond // The interval at which the simulated blockchain ticks.
const simChainTick = 15 * time.Second    // The amount of time that is added each blockchain tick.

func TestPaymentHappy(t *testing.T) {
	rng := pkgtest.Prng(t)

	names := [2]string{"Alice", "Bob"}
	asset := pchtest.NewRandomAsset(rng)
	amounts := [2]*big.Int{big.NewInt(100), big.NewInt(100)}
	setups := newRoleSetups(t, rng, names, asset, amounts[:])

	roles := [2]ctest.Executer{
		ctest.NewAlice(setups[0], t),
		ctest.NewBob(setups[1], t),
	}

	peers := [2]wire.Address{
		setups[0].Identity.Address(),
		setups[1].Identity.Address(),
	}

	cfg := ctest.AliceBobExecConfig{
		BaseExecConfig: ctest.MakeBaseExecConfig(
			peers,
			asset,
			amounts,
			client.WithoutApp(),
		),
		NumPayments: [2]int{2, 2},
		TxAmounts:   [2]*big.Int{big.NewInt(5), big.NewInt(3)},
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()
	err := ctest.ExecuteTwoPartyTest(ctx, roles, &cfg)
	if err != nil {
		t.Error(err)
	}
}

func TestPaymentDispute(t *testing.T) {
	rng := pkgtest.Prng(t)

	names := [2]string{"Mallory", "Carol"}
	asset := pchtest.NewRandomAsset(rng)
	amounts := [2]*big.Int{big.NewInt(100), big.NewInt(100)}
	setups := newRoleSetups(t, rng, names, asset, amounts[:])

	roles := [2]ctest.Executer{
		ctest.NewMallory(setups[0], t),
		ctest.NewCarol(setups[1], t),
	}

	peers := [2]wire.Address{
		setups[0].Identity.Address(),
		setups[1].Identity.Address(),
	}

	cfg := ctest.MalloryCarolExecConfig{
		BaseExecConfig: ctest.MakeBaseExecConfig(
			peers,
			asset,
			amounts,
			client.WithoutApp(),
		),
		NumPayments: [2]int{5, 0},
		TxAmounts:   [2]*big.Int{big.NewInt(20), big.NewInt(0)},
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()
	err := ctest.ExecuteTwoPartyTest(ctx, roles, &cfg)
	if err != nil {
		t.Error(err)
	}
}

func newRoleSetups(t *testing.T, rng *rand.Rand, names [2]string, asset channel.Asset, amounts []channel.Bal) []ctest.RoleSetup {
	var (
		bus   = wiretest.NewSerializingLocalBus()
		n     = len(names)
		setup = make([]ctest.RoleSetup, n)
	)

	// Create client.
	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()
	c, contract := chtest.NewTestClientWithContract(ctx, t)

	// Start ticking.
	c.StartTicking(blockTick, simChainTick)
	t.Cleanup(c.StopTicking)

	// Create adjudicator.
	adjOpt := bchannel.AdjudicatorPollingIntervalOpt(polling)
	adj := bchannel.NewAdjudicator(c, contract, c.Account(), adjOpt)
	funderOpt := bchannel.FunderPollingIntervalOpt(polling)
	funder := bchannel.NewFunder(c, contract, c.Account(), funderOpt)

	for i := 0; i < n; i++ {
		setup[i] = ctest.RoleSetup{
			Name:              names[i],
			Identity:          wtest.NewRandomAccount(rng),
			Bus:               bus,
			Funder:            funder,
			Adjudicator:       adj,
			Wallet:            wtest.NewWallet(),
			Timeout:           roleOperationTimeout,
			ChallengeDuration: challengeDuration,
		}

		coins := binding.MakeCoins([]channel.Asset{asset}, []channel.Bal{amounts[i]})
		err := c.AddCoins(ctx, c.Account(), coins)
		if err != nil {
			panic(err)
		}
	}

	return setup
}
