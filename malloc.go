// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build !tinygo

package nottinygc

import "C"
import "unsafe"

func cmalloc(size uintptr) unsafe.Pointer {
	return C.malloc(C.ulong(size))
}

func cfree(ptr unsafe.Pointer) {
	C.free(ptr)
}
