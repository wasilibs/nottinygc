// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build tinygo && !nottinygc_proxywasm

package nottinygc

/*
#cgo LDFLAGS: --export=malloc --export=free
*/
import "C"
