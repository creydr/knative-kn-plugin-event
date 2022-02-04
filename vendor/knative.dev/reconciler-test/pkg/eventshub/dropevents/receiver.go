/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dropevents

import (
	"sync"

	"knative.dev/reconciler-test/pkg/eventshub/dropevents/dropeventsfibonacci"
	"knative.dev/reconciler-test/pkg/eventshub/dropevents/dropeventsfirst"
)

const (
	Fibonacci = "fibonacci"
	Sequence  = "sequence"
)

// Skipper represents the logic to apply to accept/reject events.
type Skipper interface {
	// Skip returns true if the message must be rejected
	Skip(counter uint64) bool
}

type noopSkip struct{}

var NoopSkipper Skipper = noopSkip{}

func (n noopSkip) Skip(uint64) bool {
	return false
}

type CounterHandler struct {
	counter uint64
	Skipper Skipper
	sync.Mutex
}

func (h *CounterHandler) Skip() bool {
	h.Lock()
	defer h.Unlock()

	h.counter++
	return h.Skipper.Skip(h.counter)
}

// count is only used for SKIP_ALGORITHM=Sequence.
func SkipperAlgorithmWithCount(algorithm string, count uint64) Skipper {
	switch algorithm {
	case Fibonacci:
		return &dropeventsfibonacci.Fibonacci{Prev: 1, Current: 1}

	case Sequence:
		return dropeventsfirst.First{N: count}

	default:
		panic("unknown algorithm: " + algorithm)
	}
}
