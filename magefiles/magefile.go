// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test runs unit tests
func Test() error {
	v, err := sh.Output("tinygo", "version")
	if err != nil {
		return fmt.Errorf("invoking tinygo: %w", err)
	}

	tags := []string{"custommalloc"}
	if strings.HasSuffix(v, "tinygo version 0.28.") {
		tags = append(tags, "nottinygc_finalizer")
	}

	if err := sh.RunV("tinygo", "test", "-gc=custom", fmt.Sprintf("-tags='%s'", strings.Join(tags, " ")), "-target=wasi", "-v", "-scheduler=none", "./..."); err != nil {
		return err
	}

	var stdout bytes.Buffer
	if _, err := sh.Exec(map[string]string{}, &stdout, io.Discard, "tinygo", "test", "-target=wasi", "-v", "-scheduler=none", "./..."); err == nil {
		return errors.New("expected tinygo test to fail without -gc=custom")
	}
	if s := stdout.String(); !strings.Contains(s, "nottinygc requires passing -gc=custom and -tags=custommalloc to TinyGo when compiling") {
		return fmt.Errorf("unexpected error message: %s", s)
	}

	if err := buildBenchExecutable(); err != nil {
		return err
	}

	if err := sh.RunV("go", "test", "./bench"); err != nil {
		return err
	}

	return nil
}

func Format() error {
	if err := sh.RunV("go", "run", fmt.Sprintf("mvdan.cc/gofumpt@%s", verGoFumpt), "-l", "-w", "."); err != nil {
		return err
	}

	// addlicense strangely logs skipped files to stderr despite not being erroneous, so use the long sh.Exec form to
	// discard stderr too.
	if _, err := sh.Exec(map[string]string{}, io.Discard, io.Discard, "go", "run", fmt.Sprintf("github.com/google/addlicense@%s", verAddLicense),
		"-c", "wasilibs authors",
		"-l", "mit",
		"-s=only",
		"-y=",
		"-ignore", "**/*.yml",
		"-ignore", "**/*.yaml",
		"."); err != nil {
		return err
	}

	if err := sh.RunV("go", "run", fmt.Sprintf("github.com/rinchsan/gosimports/cmd/gosimports@%s", verGosImports), "-w",
		"-local", "github.com/wasilibs/nottinygc",
		"."); err != nil {
		return nil
	}
	return nil
}

func Lint() error {
	if _, err := sh.Exec(map[string]string{}, io.Discard, io.Discard, "go", "run", fmt.Sprintf("github.com/google/addlicense@%s", verAddLicense),
		"-check",
		"-c", "wasilibs authors",
		"-s=only",
		"-l=mit",
		"-y=",
		"-ignore", "**/*.yml",
		"-ignore", "**/*.yaml",
		"."); err != nil {
		return fmt.Errorf("missing copyright headers, use go run mage.go format: %w", err)
	}

	return sh.RunV("go", "run", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", verGolancCILint), "run")
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
	if err := buildBenchExecutable(); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-scheduler=none", "-target=wasi", "-o", "build/benchref.wasm", "./bench"); err != nil {
		return err
	}

	return sh.RunV("go", "test", "-bench=.", "-benchtime=10s", "./bench")
}

func buildBenchExecutable() error {
	if err := os.MkdirAll("build", 0o755); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-gc=custom", "-tags=custommalloc", "-scheduler=none", "-target=wasi", "-o", "build/bench.wasm", "./bench"); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-gc=custom", "-tags='custommalloc nottinygc_envoy'", "-scheduler=none", "-target=wasi", "-o", "build/bench_envoy.wasm", "./bench"); err != nil {
		return err
	}

	return nil
}

var Default = Test
