package factory

import (
	"context"

	"go-sample/internal/app"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// go get -u github.com/prometheus/client_golang

// NewPrometheusRegistry - comment func.
func NewPrometheusRegistry(_ context.Context, _ app.Options) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return registry
}
