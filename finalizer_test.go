// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

import (
	"runtime"
	"testing"
)

func TestFinalizer(t *testing.T) {
	finalized := 0
	allocFinalized(10, func() {
		finalized++
	})

	runtime.GC()

	if finalized == 0 {
		t.Errorf("finalizer not called")
	}
}

//go:noinline
func allocFinalized(num int, cb func()) {
	f := &finalized{
		a: 100,
		b: "foo",
	}

	runtime.SetFinalizer(f, func(f interface{}) {
		cb()
	})

	// Recurse to create some more stack frames or else the shadow stack
	// may still not have been reset, causing GC to find the inaccessible
	// pointers.
	if num > 1 {
		allocFinalized(num-1, cb)
	}
}

type finalized struct {
	a int
	b string
}
