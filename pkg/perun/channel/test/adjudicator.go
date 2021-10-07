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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

// Adjudicator represents an adjudicator test setup.
//
// Adjudicator should return the adjudicator instance used for testing.
// NewFundedChannel should return a funded channel at version 0, non-final.
// SignState should return signatures for the given state.
// Account returns the account for an address.
type Adjudicator interface {
	Adjudicator() channel.Adjudicator
	NewFundedChannel(context.Context, *rand.Rand) (channel.Params, channel.State)
	SignState(state *channel.State, parts []wallet.Address) []wallet.Sig
	Account(wallet.Address) wallet.Account
}

// TestAdjudicatorWithSubscription tests an adjudicator with its event subscription.
func TestAdjudicatorWithSubscription(ctx context.Context, t *testing.T, rng *rand.Rand, a Adjudicator) {
	TestAdjudicatorWithSubscriptionCollaborative(ctx, t, a, rng)
	TestAdjudicatorWithSubscriptionDispute(ctx, t, a, rng)
}

// TestAdjudicatorWithSubscriptionCollaborative tests an adjudicator with its event subscription.
func TestAdjudicatorWithSubscriptionCollaborative(ctx context.Context, t *testing.T, a Adjudicator, rng *rand.Rand) {
	adj := a.Adjudicator()

	// Set up new funded channel.
	params, state := a.NewFundedChannel(ctx, rng)
	expectedVersion := uint64(0)
	if state.Version != expectedVersion {
		t.Fatalf("wrong version: expected %d, got %d", expectedVersion, state.Version)
	}
	subChannels := []channel.SignedState{} // We do not test with subchannels yet.

	// Adjudicator: Subscribe to events.
	eventSub, err := adj.Subscribe(ctx, params.ID())
	require.NoError(t, err, "subscribe")

	// Adjudicator: Withdraw.
	state.IsFinal = true
	sigs := a.SignState(&state, params.Parts)
	req := channel.AdjudicatorReq{
		Params: &params,
		Tx: channel.Transaction{
			State: &state,
			Sigs:  sigs,
		},
	}
	numParts := len(params.Parts)
	for _, i := range rand.Perm(numParts) {
		req.Idx = channel.Index(i)
		req.Acc = a.Account(params.Parts[i])
		err = adj.Withdraw(ctx, req, MakeStateMapFromSignedStates(subChannels...))
		require.NoErrorf(t, err, "withdraw: part %d", i)
	}

	// Subscription: Check concluded event.
	{
		e, ok := eventSub.Next().(*channel.ConcludedEvent)
		assert.True(t, ok, "concluded")
		assert.True(t, e.ID() == params.ID(), "equal ID")
		assert.True(t, e.Version() == state.Version, "version")
		err = e.Timeout().Wait(ctx)
		assert.NoError(t, err, "concluded: wait")
	}

	// Subscription: Close.
	{
		err := eventSub.Close()
		assert.NoError(t, err, "close")
		err = eventSub.Err()
		assert.NoError(t, err, "err")
	}
}

// TestAdjudicatorWithSubscriptionDispute tests an adjudicator with its event subscription.
func TestAdjudicatorWithSubscriptionDispute(ctx context.Context, t *testing.T, a Adjudicator, rng *rand.Rand) {
	adj := a.Adjudicator()

	// Set up new funded channel.
	params, state := a.NewFundedChannel(ctx, rng)
	expectedVersion := uint64(0)
	if state.Version != expectedVersion {
		t.Fatalf("wrong version: expected %d, got %d", expectedVersion, state.Version)
	}
	subChannels := []channel.SignedState{} // We do not test with subchannels yet.

	// Adjudicator: Subscribe to events.
	eventSub, err := adj.Subscribe(ctx, params.ID())
	require.NoError(t, err, "subscribe")

	// Adjudicator: Register version 0.
	{
		sigs := a.SignState(&state, params.Parts)
		req := channel.AdjudicatorReq{
			Params: &params,
			Tx: channel.Transaction{
				State: &state,
				Sigs:  sigs,
			},
		}
		err = adj.Register(ctx, req, subChannels)
		require.NoError(t, err, "register version 0")
	}

	// Subscription: Check registered event version 0.
	{
		e, ok := eventSub.Next().(*channel.RegisteredEvent)
		assert.True(t, ok, "registered")
		assert.True(t, e.ID() == params.ID(), "equal ID")
		assert.True(t, e.Version() == state.Version, "version")
		assert.True(t, e.State.Equal(&state) == nil, "equal state")
	}

	// Adjudicator: Register version 1.
	state.Version = 1
	{
		sigs := a.SignState(&state, params.Parts)
		req := channel.AdjudicatorReq{
			Params: &params,
			Tx: channel.Transaction{
				State: &state,
				Sigs:  sigs,
			},
		}
		err = adj.Register(ctx, req, subChannels)
		require.NoError(t, err, "register version 1")
	}

	// Subscription: Check registered event version 1 and wait for timeout.
	{
		e, ok := eventSub.Next().(*channel.RegisteredEvent)
		assert.True(t, ok, "registered")
		assert.True(t, e.ID() == params.ID(), "equal ID")
		assert.True(t, e.Version() == state.Version, "version")
		assert.True(t, e.State.Equal(&state) == nil, "equal state")
		err = e.Timeout().Wait(ctx)
		assert.NoError(t, err, "registered: wait")
	}

	// Adjudicator: Progress.
	// We do not test on-chain progression yet.

	// Adjudicator: Withdraw.
	{
		sigs := a.SignState(&state, params.Parts)
		req := channel.AdjudicatorReq{
			Params: &params,
			Tx: channel.Transaction{
				State: &state,
				Sigs:  sigs,
			},
		}
		numParts := len(params.Parts)
		for _, i := range rand.Perm(numParts) {
			req.Idx = channel.Index(i)
			req.Acc = a.Account(params.Parts[i])
			err = adj.Withdraw(ctx, req, MakeStateMapFromSignedStates(subChannels...))
			require.NoErrorf(t, err, "withdraw: part %d", i)
		}
	}

	// Subscription: Check concluded event.
	{
		e, ok := eventSub.Next().(*channel.ConcludedEvent)
		assert.True(t, ok, "concluded")
		assert.True(t, e.ID() == params.ID(), "equal ID")
		assert.True(t, e.Version() == state.Version, "version")
		err = e.Timeout().Wait(ctx)
		assert.NoError(t, err, "concluded: wait")
	}

	// Subscription: Close.
	{
		err := eventSub.Close()
		assert.NoError(t, err, "close")
		err = eventSub.Err()
		assert.NoError(t, err, "err")
	}
}

func MakeStateMapFromSignedStates(channels ...channel.SignedState) channel.StateMap {
	m := channel.MakeStateMap()
	for _, c := range channels {
		m.Add(c.State)
	}
	return m
}
