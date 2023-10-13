// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build gc.custom

package nottinygc

import (
	"runtime"
	"unsafe"
)

/*
#include <stddef.h>

void* GC_malloc(unsigned int size);
void* GC_malloc_explicitly_typed(unsigned int size, void* gc_descr);
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
}

// alloc tries to find some free space on the heap, possibly doing a garbage
// collection cycle if needed. If no space is free, it panics.
//
//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	// For now, provide type information when there are no pointers. In the future,
	// we can try to make this more precise for all types but this should still handle
	// the common case of large arrays.
	layout := uintptr(layoutPtr)
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
		if layoutSz == 1 && layoutBm == 0 {
			// No pointers!
			buf := C.GC_malloc_explicitly_typed(C.uint(size), nil)
			if buf == nil {
				panic("out of memory")
			}
			return buf
		}
	}
	buf := C.GC_malloc(C.uint(size))
	if buf == nil {
		panic("out of memory")
	}
	return buf
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
