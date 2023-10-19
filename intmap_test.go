// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

import "testing"

func TestIntMapBasic(t *testing.T) {
	m := newIntMap()
	_, ok := m.get(5)
	if ok {
		t.Fatal("expected not ok in empty map")
	}

	m.put(5, 10)
	v, ok := m.get(5)
	if !ok {
		t.Fatal("expected ok in map")
	}
	if v != 10 {
		t.Fatal("expected 10 in map")
	}
}

func TestIntMapNoResize(t *testing.T) {
	top := int(0.75 * 512)
	m := newIntMap()
	for i := 0; i < top; i++ {
		m.put(uintptr(i), uintptr(i))
	}

	if len(m.buckets) != 512 {
		t.Fatal("expected 512 buckets")
	}

	for i := 0; i < top; i++ {
		v, ok := m.get(uintptr(i))
		if !ok {
			t.Fatalf("expected %d to be in map", i)
		}
		if v != uintptr(i) {
			t.Fatalf("expected %d to have value %d in map, got %d", i, i, v)
		}
	}
	_, ok := m.get(uintptr(top))
	if ok {
		t.Fatal("expected not ok in map")
	}
}

func TestIntMapResize(t *testing.T) {
	top := 512
	m := newIntMap()
	for i := 0; i < top; i++ {
		m.put(uintptr(i), uintptr(i))
	}

	if len(m.buckets) != 1024 {
		t.Fatal("expected 1024 buckets")
	}

	for i := 0; i < top; i++ {
		v, ok := m.get(uintptr(i))
		if !ok {
			t.Fatalf("expected %d to be in map", i)
		}
		if v != uintptr(i) {
			t.Fatalf("expected %d to have value %d in map, got %d", i, i, v)
		}
	}
	_, ok := m.get(uintptr(top))
	if ok {
		t.Fatal("expected not ok in map")
	}
}
