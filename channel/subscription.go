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

package channel

import (
	"context"
	"sync"
	"time"

	"github.com/perun-network/perun-cosmwasm-backend/channel/binding"
	"perun.network/go-perun/channel"
)

// EventSubscription provides methods for consuming channel events.
type EventSubscription struct {
	adjudicator *Adjudicator
	channelID   channel.ID
	closed      chan struct{}
	err         chan error
	prev        binding.DisputeQueryResponse
	once        sync.Once
}

func NewEventSubscription(a *Adjudicator, ch channel.ID) *EventSubscription {
	return &EventSubscription{
		adjudicator: a,
		channelID:   ch,
		closed:      make(chan struct{}),
		err:         make(chan error, 1),
		prev:        binding.DisputeQueryResponse{},
	}
}

// Next returns the most recent or next future event. If the subscription is
// closed or any other error occurs, it returns nil.
func (s *EventSubscription) Next() channel.AdjudicatorEvent {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan := make(chan channel.AdjudicatorEvent)
	errChan := make(chan error)

	go func() {
		for {
			d, err := s.readState(ctx)
			if err != nil && err.Error() != "Unknown dispute: query wasm contract failed" {
				errChan <- err
				return
			}

			if !d.Equal(s.prev) {
				s.prev = d
				eventChan <- s.makeEvent(d)
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(s.adjudicator.polling):
			}
		}
	}()

	select {
	case <-s.closed:
		return nil
	case err := <-errChan:
		s.err <- err
		return nil
	case e := <-eventChan:
		return e
	}
}

// Err returns the error status of the subscription. After Next returns nil,
// Err should be checked for an error.
func (s *EventSubscription) Err() error {
	select {
	case <-s.closed:
		return nil
	default:
		return <-s.err
	}
}

// Close closes the subscription.
func (s *EventSubscription) Close() error {
	s.once.Do(func() { close(s.closed) })
	return nil
}

func (s *EventSubscription) readState(ctx context.Context) (binding.DisputeQueryResponse, error) {
	ch := s.channelID
	q, err := binding.NewDisputeQueryMsg(ch)
	if err != nil {
		return binding.DisputeQueryResponse{}, err
	}

	resp, err := s.adjudicator.Query(ctx, q)
	if err != nil {
		return binding.DisputeQueryResponse{}, err
	}

	return binding.DecodeDisputeQueryResponse(resp.Data)
}

func (s *EventSubscription) makeEvent(d binding.DisputeQueryResponse) channel.AdjudicatorEvent {
	state := d.State.PerunState()
	cID := state.ID
	v := state.Version
	timeout := makeTimeout(s.adjudicator.client, d.Timeout(), s.adjudicator.polling)
	if d.Concluded {
		return channel.NewConcludedEvent(cID, timeout, v)
	}
	return channel.NewRegisteredEvent(cID, timeout, v, state, nil)
}
