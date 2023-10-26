// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

package internal

import (
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

type Metrics struct {
	counters map[string]proxywasm.MetricCounter
}

const MetricPrefix = "envoy_wasm_plugin"

func NewMetrics() *Metrics {
	return &Metrics{
		counters: make(map[string]proxywasm.MetricCounter),
	}
}

func (m *Metrics) Increment(name string, tags [][2]string) {
	fullName := metricName(name, tags)
	if _, exists := m.counters[fullName]; !exists {
		m.counters[fullName] = proxywasm.DefineCounterMetric(fullName)
	}
	m.counters[fullName].Increment(1)
}

func metricName(name string, tags [][2]string) string {
	fullName := fmt.Sprintf("%s_%s", MetricPrefix, name)

	for _, t := range tags {
		fullName += fmt.Sprintf("_%s=.=%s;.;", t[0], t[1])
	}
	return fullName
}
