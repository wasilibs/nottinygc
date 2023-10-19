// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package nottinygc

import "unsafe"

// CPP_WORDSZ is a simple integer constant representing the word size
const cppWordsz = unsafe.Sizeof(uintptr(0)) * 8

type gcBitmap struct {
	words []uintptr
}

func newBitmap(size uintptr) gcBitmap {
	bmSize := gcBitmapSize(size)
	wordsArr := cmalloc(bmSize * unsafe.Sizeof(uintptr(0)))
	words := unsafe.Slice((*uintptr)(wordsArr), bmSize)
	for i := 0; i < len(words); i++ {
		words[i] = 0
	}
	return gcBitmap{words: words}
}

func (b gcBitmap) set(idx uintptr) {
	b.words[idx/cppWordsz] |= 1 << (idx % cppWordsz)
}

func (b gcBitmap) get(idx uintptr) uintptr {
	return (b.words[idx/cppWordsz] >> (idx % cppWordsz)) & 1
}

func gcBitmapSize(size uintptr) uintptr {
	return (size + cppWordsz - 1) / cppWordsz
}
