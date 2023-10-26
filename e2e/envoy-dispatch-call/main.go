// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	_ "github.com/wasilibs/nottinygc"
	"github.com/wasilibs/nottinygc/e2e/envoy-dispatch-call/internal"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext Override types.DefaultVMContext otherwise this plugin would do nothing :)
func (v *vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &filterContext{metrics: internal.NewMetrics()}
}

type filterContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext

	conf    *internal.Config
	metrics *internal.Metrics
}

// OnPluginStart Override types.DefaultPluginContext.
func (h *filterContext) OnPluginStart(_ int) types.OnPluginStartStatus {
	h.conf = internal.NewConfig()
	return types.OnPluginStartStatusOK
}

// NewHttpContext Override types.DefaultPluginContext to allow us to declare a request handler for each
// intercepted request the Envoy Sidecar sends us
func (h *filterContext) NewHttpContext(_ uint32) types.HttpContext {
	return &internal.RequestHandler{Conf: h.conf, Metrics: h.metrics}
}
