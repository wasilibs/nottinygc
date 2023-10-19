// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

import "C"
import "unsafe"

// Aim for initial over on the order of a few kilobytes.

const initialBuckets = 512
const numEmbedded = 8

type item struct {
	key uintptr
	val uintptr
}

type extraNode struct {
	next *extraNode
	item item
}

type bucket struct {
	embedded [numEmbedded]item
	extra    *extraNode
	count    byte
}

// intMap is a map from int to int. As it is used to cache descriptors within
// allocation, it cannot itself allocate using the Go heap and uses malloc
// instead. It also takes advantage of knowing we never replace values, so it
// does not support replacement.
type intMap struct {
	buckets []bucket
	count   int
}

func newIntMap() intMap {
	return intMap{
		buckets: newBuckets(initialBuckets),
		count:   0,
	}
}

func (m *intMap) put(key uintptr, val uintptr) {
	if float64(m.count+1) > float64(len(m.buckets))*0.75 {
		m.resize()
	}
	doPut(m.buckets, key, val)
	m.count++
}

func doPut(buckets []bucket, key uintptr, val uintptr) {
	pos := hash(key) % uintptr(len(buckets))
	b := &buckets[pos]
	if b.count < numEmbedded {
		b.embedded[b.count] = item{key: key, val: val}
	} else {
		e := newExtraNode()
		e.item = item{key: key, val: val}
		e.next = b.extra
		b.extra = e
	}
	b.count++
}

func (m *intMap) resize() {
	newSz := len(m.buckets) * 2
	newBkts := newBuckets(newSz)
	for i := 0; i < len(m.buckets); i++ {
		b := &m.buckets[i]
		for j := 0; j < int(b.count); j++ {
			if j < numEmbedded {
				doPut(newBkts, b.embedded[j].key, b.embedded[j].val)
			} else {
				for n := b.extra; n != nil; {
					doPut(newBkts, n.item.key, n.item.val)
					next := n.next
					cfree(unsafe.Pointer(n))
					n = next
				}
			}
		}
	}
	cfree(unsafe.Pointer(&m.buckets[0]))
	m.buckets = newBkts
}

func (m *intMap) get(key uintptr) (uintptr, bool) {
	pos := hash(key) % uintptr(len(m.buckets))
	b := &m.buckets[pos]
	for i := 0; i < int(b.count); i++ {
		if i < numEmbedded {
			if b.embedded[i].key == key {
				return b.embedded[i].val, true
			}
		} else {
			for n := b.extra; n != nil; n = n.next {
				if n.item.key == key {
					return n.item.val, true
				}
			}
			break
		}
	}
	return 0, false
}

func hash(key uintptr) uintptr {
	// Use Java algorithm for cheap and easy handling of aligned values.
	// There are better ones with more operations, we can try later if needed.
	return key ^ (key >> 16)
}

func newBuckets(size int) []bucket {
	sz := unsafe.Sizeof(bucket{}) * uintptr(size)
	bucketsArr := cmalloc(sz)
	for i := uintptr(0); i < sz; i++ {
		*(*byte)(unsafe.Pointer(uintptr(bucketsArr) + i)) = 0
	}
	buckets := unsafe.Slice((*bucket)(bucketsArr), size)
	return buckets
}

func newExtraNode() *extraNode {
	sz := unsafe.Sizeof(extraNode{})
	arr := cmalloc(sz)
	for i := uintptr(0); i < sz; i++ {
		*(*byte)(unsafe.Pointer(uintptr(arr) + i)) = 0
	}
	return (*extraNode)(arr)
}
