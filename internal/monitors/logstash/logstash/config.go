package logstash

import (
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-agent/internal/core/config"
)

// Config for this monitor
type Config struct {
	config.MonitorConfig `yaml:",inline" acceptsEndpoints:"true" singleInstance:"false"`
	// The hostname of Logstash monitoring API
	Host string `yaml:"host" default:"127.0.0.1"`
	// The port number of Logstash monitoring API
	Port uint16 `yaml:"port" default:"9600"`
	// If true, the agent will connect to the host using HTTPS instead of plain HTTP.
	UseHTTPS bool `yaml:"useHTTPS"`
	// The maximum amount of time to wait for API requests
	TimeoutSeconds int `yaml:"timeoutSeconds" default:"5"`
	// Whether it will send all extra hot threads metrics as well.
	EnableExtraHotThreadsMetrics bool `yaml:"enableExtraHotThreadsMetrics"`
	// Whether it will send all extra JVM metrics as well.
	EnableExtraJVMMetrics bool `yaml:"enableExtraJVMMetrics"`
	// Whether it will send all extra OS metrics as well.
	EnableExtraOSMetrics bool `yaml:"enableExtraOSMetrics"`
	// Whether it will send all extra pipelines metrics as well.
	EnableExtraPipelinesMetrics bool `yaml:"enableExtraPipelinesMetrics"`
	// Whether it will send all extra events metrics as well.
	EnableExtraEventsMetrics bool `yaml:"enableExtraEventsMetrics"`
	// Whether it will send all extra process metrics as well.
	EnableExtraProcessMetrics bool `yaml:"enableExtraProcessMetrics"`
	// Whether it will send all extra reloads metrics as well.
	EnableExtraReloadsMetrics bool `yaml:"enableExtraReloadsMetrics"`
}

// GetExtraMetrics handles the legacy enhancedMetrics option.
func (c *Config) GetExtraMetrics() []string {
	var extraMetrics []string

	enabledGroups := c.getEnabledGroups()

	for groupName, enabled := range enabledGroups {
		if enabled {
			extraMetrics = append(extraMetrics, groupMetricsMap[groupName]...)
		}
	}

	return extraMetrics
}

var _ config.ExtraMetrics = &Config{}

func (c *Config) getMetricTypeMap() map[string]datapoint.MetricType {
	metricTypeMap := make(map[string]datapoint.MetricType)

	for metricName := range defaultMetrics {
		metricTypeMap[metricName] = metricSet[metricName].Type
	}

	enabledGroups := c.getEnabledGroups()

	for groupName, enabled := range enabledGroups {
		if enabled {
			for _, metricName := range groupMetricsMap[groupName] {
				metricTypeMap[metricName] = metricSet[metricName].Type
			}
		}
	}

	return metricTypeMap
}

func (c *Config) getEnabledGroups() map[string]bool {
	return map[string]bool{
		groupHotThreads: c.EnableExtraHotThreadsMetrics,
		groupJvm:        c.EnableExtraJVMMetrics,
		groupOs:         c.EnableExtraOSMetrics,
		groupPipeline:   c.EnableExtraPipelinesMetrics,
		groupEvents:     c.EnableExtraEventsMetrics,
		groupProcess:    c.EnableExtraProcessMetrics,
		groupReloads:    c.EnableExtraReloadsMetrics,
	}
}
