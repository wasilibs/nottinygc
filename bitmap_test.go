// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

import (
	"fmt"
	"testing"
)

func TestBitmap32Bits(t *testing.T) {
	tests := []uintptr{
		0,
		0b1,
		0b101,
		0b111,
		0b0001,
		0b1000001,
		0xFFFFFFFF,
		0x11111111,
		0x01010101,
		0x0F0F0F0F,
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			bm := newBitmap(32)
			if len(bm.words) != 1 {
				t.Fatalf("expected 1 word, got %v", len(bm.words))
			}
			for i := 0; i < 32; i++ {
				if tc&(1<<i) != 0 {
					bm.set(uintptr(i))
				}
			}

			for i := 0; i < 32; i++ {
				got := bm.get(uintptr(i))
				if tc&(1<<i) != 0 {
					if got == 0 {
						t.Fatalf("expected bit %v to be set", i)
					}
				} else {
					if got != 0 {
						t.Fatalf("expected bit %v to be unset", i)
					}
				}
			}
		})
	}
}

// Test for multiple words, we pick larger than 64-bits to have more than one word on Go
// as well. We don't actually run CI with Go but it can be helpful for development.
func TestBitmap128Bits(t *testing.T) {
	// We'll just repeat these.
	tests := []uintptr{
		0,
		0b1,
		0b101,
		0b111,
		0b0001,
		0b1000001,
		0xFFFFFFFF,
		0x11111111,
		0x01010101,
		0x0F0F0F0F,
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			bm := newBitmap(128)
			if cppWordsz == 32 && len(bm.words) != 4 || cppWordsz == 64 && len(bm.words) != 2 {
				t.Fatalf("got %v words", len(bm.words))
			}
			for j := 0; j < 4; j++ {
				for i := 0; i < 32; i++ {
					if tc&(1<<(32*j+i)) != 0 {
						bm.set(uintptr(i))
					}
				}
			}

			for j := 0; j < 4; j++ {
				for i := 0; i < 32; i++ {
					got := bm.get(uintptr(32*j + i))
					if tc&(1<<(32*j+i)) != 0 {
						if got == 0 {
							t.Fatalf("expected bit %v to be set", i)
						}
					} else {
						if got != 0 {
							t.Fatalf("expected bit %v to be unset", i)
						}
					}
				}
			}
		})
	}
}
