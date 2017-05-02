package config

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

// Plugin describes a collectd plugin
type Plugin struct {
	Plugin       string
	Name         string
	Template     string
	TemplateFile string
	Dims         string
	Host         string
	Port         uint16
	Vars         map[string]interface{}
}

// PluginType is a type for collectd plugins
type PluginType string

// Define static collectd plugins
const (
	// SignalFx plugin
	SignalFx PluginType = "signalfx"
	// WriteHTTP plugin
	WriteHTTP PluginType = "writehttp"
	// Docker plugin
	Docker PluginType = "docker"
	// MesosMaster plugin
	MesosMaster PluginType = "mesos-master"
	// MesosAgent plugin
	MesosAgent PluginType = "mesos-agent"
	// Marathon plugin
	Marathon PluginType = "marathon"
)

// PLUGINS is a mapping to create plugin instances with defaults
var PLUGINS = map[PluginType]func(string) *Plugin{
	SignalFx: func(instanceName string) *Plugin {
		// XXX: Super hacky. Ideally this should have no knowledge of the global
		// viper.
		return &Plugin{
			Plugin:       "signalfx",
			TemplateFile: "signalfx.conf.tmpl",
			Vars: map[string]interface{}{
				"url": viper.GetString("ingesturl"),
			},
			Name: instanceName,
		}
	},
	WriteHTTP: func(instanceName string) *Plugin {
		// XXX: Super hacky. Ideally this should have no knowledge of the global
		// viper.
		query := url.Values{}

		for k, v := range viper.GetStringMapString("dimensions") {
			query["sfxdim_"+k] = []string{v}
		}

		plugin := &Plugin{
			Plugin:       "write-http",
			TemplateFile: "write-http.conf.tmpl",
			Vars: map[string]interface{}{
				"url":        viper.GetString("ingesturl"),
				"dimensions": query.Encode(),
			},
			Name: instanceName,
		}

		return plugin
	},
	Docker: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:       "docker",
			TemplateFile: "docker.conf.tmpl",
			Name:         instanceName,
			Vars: map[string]interface{}{
				"url": "unix:///var/run/docker.sock",
			},
		}
	},
	MesosMaster: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:       "mesos-master",
			TemplateFile: "mesos-master.conf.tmpl",
			Name:         instanceName,
			Vars: map[string]interface{}{
				"host":         "localhost",
				"port":         "5050",
				"cluster":      "cluster-0",
				"instance":     "master-0",
				"systemHealth": "false",
				"verbose":      "false",
			},
		}
	},
	MesosAgent: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:       "mesos-agent",
			TemplateFile: "mesos-agent.conf.tmpl",
			Name:         instanceName,
			Vars: map[string]interface{}{
				"host":     "localhost",
				"port":     "5051",
				"cluster":  "cluster-0",
				"instance": "agent-0",
				"verbose":  "false",
			},
		}
	},
	Marathon: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:       "marathon",
			TemplateFile: "marathon.conf.tmpl",
			Name:         instanceName,
			Vars: map[string]interface{}{
				"host":     "localhost",
				"port":     "8080",
				"username": "",
				"password": "",
			},
		}
	},
}

// NewPlugin constructs a plugin with default values depending on the service type
func NewPlugin(pluginType PluginType, pluginName string) (*Plugin, error) {
	if create, ok := PLUGINS[pluginType]; ok {
		return create(pluginName), nil
	}
	return nil, fmt.Errorf("plugin %s is unsupported", pluginType)
}

// NewInstancePlugin creates a plugin for a supported service type
func NewInstancePlugin(pluginType string, pluginName string) (*Plugin, error) {
	// TODO: Maintain a list of supported service types for collectd if not all monitors support the same ones.
	return &Plugin{Plugin: pluginType, Name: pluginName, Vars: map[string]interface{}{}}, nil
}

// GroupByPlugin creates a map of instances by plugin
func GroupByPlugin(instances []*Plugin) map[string][]*Plugin {
	pluginMap := map[string][]*Plugin{}
	for _, instance := range instances {
		plugin := instance.Plugin

		if val, ok := pluginMap[plugin]; ok {
			pluginMap[plugin] = append(val, instance)
		} else {
			pluginMap[plugin] = []*Plugin{instance}
		}
	}
	return pluginMap
}

// CollectdConfig are global collectd settings
type CollectdConfig struct {
	Interval             uint
	Timeout              uint
	ReadThreads          uint
	WriteQueueLimitHigh  uint `yaml:"writeQueueLimitHigh"`
	WriteQueueLimitLow   uint `yaml:"writeQueueLimitLow"`
	CollectInternalStats bool
	Hostname             string
	Plugins              []map[string]interface{}
}

// AppConfig is the top-level configuration object consumed by templates.
type AppConfig struct {
	AgentConfig *CollectdConfig
	Plugins     map[string][]*Plugin
}

// NewCollectdConfig creates a default collectd config instance
func NewCollectdConfig() *CollectdConfig {
	return &CollectdConfig{
		Interval:             15,
		Timeout:              2,
		ReadThreads:          5,
		WriteQueueLimitHigh:  500000,
		WriteQueueLimitLow:   400000,
		CollectInternalStats: true,
	}
}
