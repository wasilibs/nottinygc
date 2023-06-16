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

// TestExportsMalloc is in the bench package since it already provides a fully packaged binary.
// This could be cleaned up by using a different package, but for now we'll stick with being a little
// lazy.
func TestExportsMalloc(t *testing.T) {
	wasm, err := os.ReadFile(filepath.Join("..", "build", "bench.wasm"))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.InstantiateModuleFromBinary(ctx, wasm)
	if err != nil {
		t.Fatal(err)
	}

	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	if malloc == nil {
		t.Error("malloc is not exported")
	}
	if free == nil {
		t.Error("free is not exported")
	}
}

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
