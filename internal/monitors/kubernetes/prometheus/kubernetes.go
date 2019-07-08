// Implements a ConfigInterface that other monitors can extend to gather Prometheus metrics using Kubernetes
// service account.
package prometheus

import (
	"fmt"
	"io"
	"time"

	k8s "k8s.io/client-go/kubernetes"

	"github.com/signalfx/signalfx-agent/internal/core/common/kubernetes"
	"github.com/signalfx/signalfx-agent/internal/core/config"

	"github.com/prometheus/common/expfmt"
	"github.com/signalfx/signalfx-agent/internal/monitors/prometheusexporter"
)

// Config for this monitor
type Config struct {
	config.MonitorConfig `yaml:",inline" acceptsEndpoints:"true"`

	// Host of the exporter
	Host string `yaml:"host" validate:"required"`
	// Port of the exporter
	Port uint16 `yaml:"port" validate:"required"`

	// Configuration of the Kubernetes API client.
	KubernetesAPI *kubernetes.APIConfig `yaml:"kubernetesAPI" default:"{}"`
	// Path to the metrics endpoint on server, usually `/metrics` (the default).
	MetricPath string `yaml:"metricPath" default:"/metrics"`

	// If true, the agent will connect to the exporter using HTTPS instead of
	// plain HTTP.
	UseHTTPS bool `yaml:"useHTTPS"`

	monitorType string
}

// Validate k8s-specific configuration.
func (c *Config) Validate() error {
	return c.KubernetesAPI.Validate()
}

// NewConfig creates an instance for a specific monitor
func NewConfig(monitorType string) *Config {
	return &Config{monitorType: monitorType}
}

// NewClient is a ConfigInterface method implementation that creates the prometheus client.
func (c *Config) NewClient() (*prometheusexporter.Client, error) {
	conf, err := kubernetes.CreateRestConfig(c.KubernetesAPI)
	if err != nil {
		return nil, err
	}

	scheme := "http"
	if c.UseHTTPS {
		scheme = "https"
	}

	conf.Host = fmt.Sprintf("%s://%s:%d", scheme, c.Host, c.Port)

	k8sClient, err := k8s.NewForConfig(conf)
	if err != nil {
		return nil, err
	}
	return &prometheusexporter.Client{
		GetBodyReader: func() (bodyReader io.ReadCloser, format expfmt.Format, err error) {
			format = expfmt.FmtText
			bodyReader, err = k8sClient.RESTClient().Get().AbsPath(c.MetricPath).Stream()
			return
		},
	}, nil
}

// GetInterval is a ConfigInterface method implementation for getting the configured monitor run interval.
func (c *Config) GetInterval() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}

// GetMonitorType is a ConfigInterface method implementation for getting the monitor type.
func (c *Config) GetMonitorType() string {
	return c.monitorType
}
