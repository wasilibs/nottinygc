// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/magefile/mage/sh"
)

func E2eCoraza() error {
	if _, err := os.Stat(filepath.Join("e2e", "coraza-proxy-wasm")); os.IsNotExist(err) {
		// Try not pinning version, there should be no compatibility issues causing unexpected failures from a
		// green coraza build so we get to keep forward coverage this way.
		if err := sh.RunV("git", "clone", "https://github.com/corazawaf/coraza-proxy-wasm.git", filepath.Join("e2e", "coraza-proxy-wasm")); err != nil {
			return err
		}
	}

	if err := os.Chdir(filepath.Join("e2e", "coraza-proxy-wasm")); err != nil {
		return err
	}
	defer func() {
		for _, f := range []string{"ftw-envoy.log"} {
			content, err := os.ReadFile(filepath.Join("build", f))
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile(filepath.Join("..", "..", "build", "logs", f), content, 0o644); err != nil {
				panic(err)
			}
		}
	}()

	if err := sh.RunV("go", "mod", "edit", "-replace=github.com/wasilibs/nottinygc=../.."); err != nil {
		return err
	}
	defer func() {
		if err := sh.RunV("go", "mod", "edit", "-dropreplace=github.com/wasilibs/nottinygc"); err != nil {
			panic(err)
		}
	}()

	if err := sh.RunV("go", "run", "mage.go", "build"); err != nil {
		return err
	}

	if err := sh.RunV("go", "run", "mage.go", "ftw"); err != nil {
		return err
	}

	return nil
}

func E2eEnvoyDispatchCall() error {
	if err := os.MkdirAll(filepath.Join("e2e", "envoy-dispatch-call", "build"), 0o755); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-target=wasi", "-gc=custom", "-tags='custommalloc nottinygc_envoy'", "-scheduler=none",
		"-o", filepath.Join("e2e", "envoy-dispatch-call", "build", "plugin.wasm"), "./e2e/envoy-dispatch-call"); err != nil {
		return err
	}

	if err := sh.RunV("docker-compose", "--file", filepath.Join("e2e", "envoy-dispatch-call", "docker-compose.yml"), "up", "-d"); err != nil {
		return err
	}
	defer func() {
		if err := sh.RunV("docker-compose", "--file", filepath.Join("e2e", "envoy-dispatch-call", "docker-compose.yml"), "down", "-v"); err != nil {
			panic(err)
		}
	}()

	stats, err := e2eLoad("http://localhost:8080/status/200", "http://localhost:8082/stats", 40, 5000)
	if err != nil {
		return err
	}

	requestCount := 0
	authCallbackCount := 0
	authSuccessCount := 0
	for _, s := range stats.Stats {
		switch s.Name {
		case "wasmcustom.envoy_wasm_plugin_on_http_request_headers_count":
			requestCount = s.Value
		case "wasmcustom.envoy_wasm_plugin_authCallback_count":
			authCallbackCount = s.Value
		case "wasmcustom.envoy_wasm_plugin_authCallback_success_count":
			authSuccessCount = s.Value
		}
	}
	if requestCount == 0 || authCallbackCount == 0 || authSuccessCount == 0 {
		return fmt.Errorf("invalid stats: %v", stats)
	}

	if authCallbackCount != requestCount {
		return fmt.Errorf("expected authCallback_count to equal request count, got %d != %d", authCallbackCount, requestCount)
	}
	if authSuccessCount != requestCount {
		return fmt.Errorf("expected authSuccess_count to equal request count, got %d != %d", authSuccessCount, requestCount)
	}

	return nil
}

func E2eHigressGCTest() error {
	if err := os.MkdirAll(filepath.Join("e2e", "higress-gc-test", "build"), 0o755); err != nil {
		return err
	}

	if err := sh.RunV("tinygo", "build", "-target=wasi", "-gc=custom", "-tags='custommalloc nottinygc_envoy'", "-scheduler=none",
		"-o", filepath.Join("e2e", "higress-gc-test", "build", "plugin.wasm"), "./e2e/higress-gc-test"); err != nil {
		return err
	}

	if err := sh.RunV("docker-compose", "--file", filepath.Join("e2e", "higress-gc-test", "docker-compose.yml"), "up", "-d"); err != nil {
		return err
	}
	defer func() {
		if err := sh.RunV("docker-compose", "--file", filepath.Join("e2e", "higress-gc-test", "docker-compose.yml"), "down", "-v"); err != nil {
			panic(err)
		}
	}()

	_, err := e2eLoad("http://localhost:8080/hello", "http://localhost:8082/stats", 4, 10000)
	if err != nil {
		return err
	}

	type memStats struct {
		Sys int `json:"Sys"`
	}

	res, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var stats memStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return err
	}

	// We expect around 20MB per VM (this reports per VM stat), a conservative
	// 100MB should be a fine check without flakiness
	if mem := stats.Sys; mem > 100_000_000 {
		return fmt.Errorf("expected <100MB memory used, actual: %d", mem)
	}

	return nil
}

type counterStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type counterStats struct {
	Stats []counterStat `json:"stats"`
}

// If needed, we can try being more sophisticated later but run some simple load for now.
func e2eLoad(url string, statsURL string, p int, n int) (*counterStats, error) {
	wg := sync.WaitGroup{}

	var success atomic.Uint32
	var fails atomic.Uint32

	healthy := false
	// Wait for healthy
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		if res.StatusCode == http.StatusOK {
			healthy = true
			break
		}
	}

	if !healthy {
		return nil, errors.New("failed to get healthy in 100 attempts")
	}

	for i := 0; i < p; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < n; j++ {
				res, err := http.Get(url)
				switch {
				case err != nil:
					fallthrough
				case res.StatusCode != http.StatusOK:
					fails.Add(1)
				default:
					success.Add(1)
				}
			}
		}()
	}

	wg.Wait()

	res, err := http.Get(fmt.Sprintf("%s?filter=wasmcustom&format=json", statsURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var stats counterStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, err
	}

	if s := success.Load(); s != uint32(p*n) {
		return &stats, fmt.Errorf("expected all requests to succeed, got success=%d, fails=%d, stats=%v", s, fails.Load(), stats)
	}

	return &stats, nil
}

func init() {
	if err := os.MkdirAll(filepath.Join("build", "logs"), 0o755); err != nil {
		panic(err)
	}
}
