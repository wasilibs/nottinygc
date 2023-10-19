// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

/*
#include <stddef.h>

void* malloc(size_t size);
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
