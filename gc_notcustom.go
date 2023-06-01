// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build !gc.custom

package nottinygc

func init() {
	panic("nottinygc requires passing -gc=custom and -tags=custommalloc to TinyGo when compiling.\nhttps://github.com/wasilibs/nottinygc#usage")
}
