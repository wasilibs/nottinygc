// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package internal

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type RequestHandler struct {
	// Bring in the callback functions
	types.DefaultHttpContext

	Conf    *Config
	Metrics *Metrics
}

const (
	XRequestIdHeader = "x-request-id"
	AuthHeader       = "authorization"
)

func (r *RequestHandler) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	r.Metrics.Increment("on_http_request_headers_count", nil)

	// None of the parameters are useful here, so we have to ask the Envoy Sidecar for the actual request headers
	requestHeaders, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		// Allow Envoy Sidecar to forward this request to the upstream service
		return types.ActionContinue
	}

	// Making this a map makes accessing specific headers much easier later on
	reqHeaderMap := headerArrayToMap(requestHeaders)

	// Grab the always-present xRequestID to help grouping logs belonging to same request
	xRequestID := reqHeaderMap[XRequestIdHeader]

	// Now we can take action on this request
	return r.doSomethingWithRequest(reqHeaderMap, xRequestID)
}

// headerArrayToMap is a simple function to convert from array of headers to a Map
func headerArrayToMap(requestHeaders [][2]string) map[string]string {
	headerMap := make(map[string]string)
	for _, header := range requestHeaders {
		headerMap[header[0]] = header[1]
	}
	return headerMap
}

func (r *RequestHandler) doSomethingWithRequest(reqHeaderMap map[string]string, xRequestID string) types.Action {
	authClient := AuthClient{XRequestID: xRequestID, Conf: r.Conf, Metrics: r.Metrics}
	authClient.RequestJWT()

	return types.ActionPause
}
