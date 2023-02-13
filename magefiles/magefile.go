package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test runs unit tests
func Test() error {
	return sh.RunV("tinygo", "test", "-gc=custom", "-tags=custommalloc", "-target=wasi", "-v", "-scheduler=none", "./...")
}

func Format() error {
	if err := sh.RunV("go", "run", fmt.Sprintf("mvdan.cc/gofumpt@%s", gofumptVersion), "-l", "-w", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "run", fmt.Sprintf("github.com/rinchsan/gosimports/cmd/gosimports@%s", gosImportsVer), "-w",
		"-local", "github.com/wasilibs/nottinygc",
		"."); err != nil {
		return nil
	}
	return nil
}

func Lint() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", golangCILintVer), "run")
}

// Check runs lint and tests.
func Check() {
	mg.SerialDeps(Lint, Test)
}

// UpdateLibs updates the precompiled wasm libraries.
func UpdateLibs() error {
	libs := []string{"bdwgc", "mimalloc"}
	for _, lib := range libs {
		if err := sh.RunV("docker", "build", "-t", "ghcr.io/wasilibs/nottinygc/buildtools-"+lib, filepath.Join("buildtools", lib)); err != nil {
			return err
		}
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := sh.RunV("docker", "run", "-it", "--rm", "-v", fmt.Sprintf("%s:/out", filepath.Join(wd, "wasm")), "ghcr.io/wasilibs/nottinygc/buildtools-"+lib); err != nil {
			return err
		}
	}
	return nil
}

// Bench runs benchmarks.
func Bench() error {
	if err := os.MkdirAll("build", 0o755); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-scheduler=none", "-target=wasi", "-o", "build/benchref.wasm", "./bench"); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-gc=custom", "-tags=custommalloc", "-scheduler=none", "-target=wasi", "-o", "build/bench.wasm", "./bench"); err != nil {
		return err
	}

	return sh.RunV("go", "test", "-bench=.", "-benchtime=10s", "./bench")
}

var Default = Test
