// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

/*
#cgo LDFLAGS: -Lwasm -lgc -lmimalloc -lclang_rt.builtins-wasm32 --export=malloc --export=free
*/
import "C"
