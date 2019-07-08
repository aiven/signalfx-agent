package scheduler

import (
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/kubernetes/prometheus"
	"github.com/signalfx/signalfx-agent/internal/monitors/prometheusexporter"
)

func init() {
	monitors.Register(
		&monitorMetadata, func() interface{} { return &prometheusexporter.Monitor{} },
		prometheus.NewConfig(monitorType))
}
