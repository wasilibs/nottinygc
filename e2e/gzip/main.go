// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"compress/gzip"

	_ "github.com/wasilibs/nottinygc"
)

func main() {
	r := bytes.NewReader([]byte{})
	gzip.NewReader(r)
}
