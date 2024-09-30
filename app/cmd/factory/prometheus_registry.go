package factory

import (
	"context"

	"github.com/mondegor/go-sample/internal/app"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// go get -u github.com/prometheus/client_golang

// NewPrometheusRegistry - создаёт объект prometheus.Registry.
func NewPrometheusRegistry(_ context.Context, opts app.Options) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	opts.InternalRouter.Handle(
		"/metrics",
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}),
	)

	return registry
}
