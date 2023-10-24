// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build tinygo && nottinygc_proxywasm

package nottinygc

/*
#cgo LDFLAGS: -Lwasm -lgc -lmimalloc -lclang_rt.builtins-wasm32
*/
import "C"

//export sched_yield
func sched_yield() int32 {
	return 0
}
