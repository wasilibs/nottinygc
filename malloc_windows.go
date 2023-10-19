// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

// TODO(anuraaga): Remove this file, it currently works around an issue with Windows but
// there should be a better way.

/*
void* malloc(unsigned long long);
void free(void* ptr);
*/
import "C"
import "unsafe"

func cmalloc(size uintptr) unsafe.Pointer {
	return C.malloc(C.ulong(size))
}

func cfree(ptr unsafe.Pointer) {
	C.free(ptr)
}
