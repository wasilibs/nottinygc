// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func BenchmarkGC(b *testing.B) {
	tests := []string{"bench.wasm", "benchref.wasm"}
	for _, tc := range tests {
		tt := tc

		wasm, err := os.ReadFile(filepath.Join("..", "build", tt))
		if err != nil {
			b.Fatal(err)
		}

		ctx := context.Background()
		r := wazero.NewRuntime(ctx)
		wasi_snapshot_preview1.MustInstantiate(ctx, r)

		mod, err := r.InstantiateModuleFromBinary(ctx, wasm)
		if err != nil {
			b.Fatal(err)
		}

		run := mod.ExportedFunction("run")

		b.Run(tt, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				run.Call(ctx)
			}
		})
	}
}
