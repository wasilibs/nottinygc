// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build !windows

package nottinygc

/*
void* malloc(unsigned int size);
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
