// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc_test

import (
	"runtime"
	"testing"
)

// Some simple GC tests copied from Go
func mk2() {
	b := new([10000]byte)
	_ = b
	//	println(b, "stored at", &b)
}

func mk1() { mk2() }

func TestGC(t *testing.T) {
	for i := 0; i < 10; i++ {
		mk1()
		runtime.GC()
	}
}

func TestGC1(t *testing.T) {
	for i := 0; i < 1e5; i++ {
		x := new([100]byte)
		_ = x
	}
}
