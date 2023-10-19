// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func cmalloc(size uintptr) unsafe.Pointer {
	return C.malloc(C.ulong(size))
}

func cfree(ptr unsafe.Pointer) {
	C.free(ptr)
}
