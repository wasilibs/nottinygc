// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build gc.custom

package nottinygc

import (
	"math/bits"
	"runtime"
	"unsafe"
)

/*
#include <stddef.h>

void* GC_malloc(unsigned int size);
void* GC_malloc_explicitly_typed(unsigned int size, unsigned int gc_descr);
void* GC_calloc_explicitly_typed(unsigned int nelements, unsigned int element_size, unsigned int gc_descr);
unsigned int GC_make_descriptor(unsigned int* bm, unsigned int len);
void GC_free(void* ptr);
void GC_gcollect();
void GC_set_on_collection_event(void* f);

size_t GC_get_gc_no();
void GC_get_heap_usage_safe(size_t* heap_size, size_t* free_bytes, size_t* unmapped_bytes, size_t* bytesSinceGC, size_t* totalBytes);
size_t GC_get_obtained_from_os_bytes();
void mi_process_info(size_t *elapsed_msecs, size_t *user_msecs, size_t *system_msecs, size_t *current_rss, size_t *peak_rss, size_t *current_commit, size_t *peak_commit, size_t *page_faults);

void onCollectionEvent();
*/
import "C"

const (
	gcEventStart = 0
)

const (
	// CPP_WORDSZ is a simple integer constant representing the word size
	cppWordsz  = uintptr(unsafe.Sizeof(uintptr(0)) * 8)
	signb      = uintptr(1) << (cppWordsz - 1)
	gcDsBitmap = uintptr(1)
)

//export onCollectionEvent
func onCollectionEvent(eventType uint32) {
	switch eventType {
	case gcEventStart:
		markStack()
	}
}

// Initialize the memory allocator.
//
//go:linkname initHeap runtime.initHeap
func initHeap() {
	C.GC_set_on_collection_event(C.onCollectionEvent)
	// We avoid overhead in calling GC_make_descriptor on every allocation by implementing
	// the bitmap computation in Go, but we need to call it at least once to initialize
	// typed GC itself.
	C.GC_make_descriptor(nil, 0)
}

// alloc tries to find some free space on the heap, possibly doing a garbage
// collection cycle if needed. If no space is free, it panics.
//
//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	var buf unsafe.Pointer

	layout := uintptr(layoutPtr)
	// We only handle the case where the layout is embedded because it is cheap to
	// transform into a descriptor for bdwgc. Larger layouts may need to allocate,
	// and we would want to cache the descriptor, but we cannot use a Go cache in this
	// function which is runtime.alloc. If we wanted to pass their information to bdwgc
	// efficiently, we would want LLVM to generate the descriptors in its format.
	// Because TinyGo uses a byte array for bitmap and bdwgc uses a word array, there are
	// likely endian issues in trying to use it directly.
	if layout&1 != 0 {
		// Layout is stored directly in the integer value.
		// Determine format of bitfields in the integer.
		const layoutBits = uint64(unsafe.Sizeof(layout) * 8)
		var sizeFieldBits uint64
		switch layoutBits { // note: this switch should be resolved at compile time
		case 16:
			sizeFieldBits = 4
		case 32:
			sizeFieldBits = 5
		case 64:
			sizeFieldBits = 6
		default:
			panic("unknown pointer size")
		}
		layoutSz := (layout >> 1) & (1<<sizeFieldBits - 1)
		layoutBm := layout >> (1 + sizeFieldBits)
		buf = allocTyped(size, layoutSz, layoutBm)
	} else {
		buf = C.GC_malloc(C.uint(size))
	}
	if buf == nil {
		panic("out of memory")
	}
	return buf
}

func allocTyped(allocSz uintptr, layoutSz uintptr, layoutBm uintptr) unsafe.Pointer {
	descr := gcDescr(layoutSz, layoutBm)

	itemSz := layoutSz * unsafe.Sizeof(uintptr(0))
	if itemSz == allocSz {
		return C.GC_malloc_explicitly_typed(C.uint(allocSz), C.uint(descr))
	}
	numItems := allocSz / itemSz
	return C.GC_calloc_explicitly_typed(C.uint(numItems), C.uint(itemSz), C.uint(descr))
}

// Reimplementation of the simple bitmap case from bdwgc
// https://github.com/ivmai/bdwgc/blob/806537be2dec4f49056cb2fe091ac7f7d78728a8/typd_mlc.c#L204
func gcDescr(layoutSz uintptr, layoutBm uintptr) uintptr {
	if layoutBm == 0 {
		return 0 // no pointers
	}

	// reversebits processes all bits but is branchless, unlike a looping version so appears
	// to perform a little better.
	return uintptr(bits.Reverse32(uint32(layoutBm))) | gcDsBitmap
}

//go:linkname free runtime.free
func free(ptr unsafe.Pointer) {
	C.GC_free(ptr)
}

//go:linkname markRoots runtime.markRoots
func markRoots(start, end uintptr) {
	// Roots are already registered in bdwgc so we have nothing to do here.
}

//go:linkname markStack runtime.markStack
func markStack()

// GC performs a garbage collection cycle.
//
//go:linkname GC runtime.GC
func GC() {
	C.GC_gcollect()
}

//go:linkname ReadMemStats runtime.ReadMemStats
func ReadMemStats(ms *runtime.MemStats) {
	var heapSize, freeBytes, unmappedBytes, bytesSinceGC, totalBytes C.size_t
	C.GC_get_heap_usage_safe(&heapSize, &freeBytes, &unmappedBytes, &bytesSinceGC, &totalBytes)

	var peakRSS C.size_t
	C.mi_process_info(nil, nil, nil, nil, &peakRSS, nil, nil, nil)

	gcOSBytes := C.GC_get_obtained_from_os_bytes()

	ms.Sys = uint64(peakRSS + gcOSBytes)
	ms.HeapSys = uint64(heapSize)
	ms.HeapIdle = uint64(freeBytes)
	ms.HeapReleased = uint64(unmappedBytes)
	ms.TotalAlloc = uint64(totalBytes)
}
