package logstash

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	"github.com/signalfx/signalfx-agent/internal/utils"
	log "github.com/sirupsen/logrus"
)

var logger = utils.NewThrottledLogger(log.WithFields(log.Fields{"monitorType": monitorType}), 30*time.Second)

const (
	nodePath = "/_node/"
)

var prefixPathMap = map[string]string{
	"node.os":            fmt.Sprintf("%s%s", nodePath, "os"),
	"node.jvm":           fmt.Sprintf("%s%s", nodePath, "jvm"),
	"node.hot_threads":   fmt.Sprintf("%s%s", nodePath, "hot_threads"),
	"node.stats.jvm":     fmt.Sprintf("%s%s", nodePath, "stats/jvm"),
	"node.stats.process": fmt.Sprintf("%s%s", nodePath, "stats/process"),
	"node.stats.events":  fmt.Sprintf("%s%s", nodePath, "stats/events"),
	"node.stats.reloads": fmt.Sprintf("%s%s", nodePath, "stats/reloads"),
	"node.stats.os":      fmt.Sprintf("%s%s", nodePath, "stats/os"),
}
var pluginPath = fmt.Sprintf("%s%s", nodePath, "plugins")
var pipelinePath = fmt.Sprintf("%s%s", nodePath, "pipelines")
var pipelineStatPath = fmt.Sprintf("%s%s", nodePath, "stats/pipelines")

var dimensionKeyMap = map[string]string{
	"threads": "thread",
	"outputs": "output",
	"inputs":  "input",
	"codecs":  "codec",
	"filters": "filter",
}

func init() {
	monitors.Register(&monitorMetadata, func() interface{} { return &Monitor{} }, &Config{})
}

// Monitor that accepts and forwards trace spans
type Monitor struct {
	Output        types.Output
	conf          *Config
	ctx           context.Context
	cancel        context.CancelFunc
	metricTypeMap map[string]datapoint.MetricType
}

// Configure the monitor and kick off volume metric syncing
func (m *Monitor) Configure(conf *Config) error {
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.conf = conf
	m.metricTypeMap = conf.getMetricTypeMap()

	client := &http.Client{
		Timeout: time.Duration(conf.TimeoutSeconds) * time.Second,
	}

	scheme := "http"
	if conf.UseHTTPS {
		scheme = "https"
	}

	dims := m.fetchNodeInfo(client, fmt.Sprintf("%s://%s:%d%s", scheme, m.conf.Host, m.conf.Port, nodePath))

	utils.RunOnInterval(m.ctx, func() {
		var dps []*datapoint.Datapoint

		for prefix, path := range prefixPathMap {
			dps = append(dps, m.fetchMetrics(client, fmt.Sprintf("%s://%s:%d%s", scheme, m.conf.Host, m.conf.Port, path), prefix, dims)...)
		}
		dps = append(dps, m.fetchPipelineMetrics(client, fmt.Sprintf("%s://%s:%d%s", scheme, m.conf.Host, m.conf.Port, pipelinePath), "node.pipelines", dims)...)
		dps = append(dps, m.fetchPipelineMetrics(client, fmt.Sprintf("%s://%s:%d%s", scheme, m.conf.Host, m.conf.Port, pipelineStatPath), "node.stats.pipelines", dims)...)
		dps = append(dps, m.fetchPluginMetrics(client, fmt.Sprintf("%s://%s:%d%s", scheme, m.conf.Host, m.conf.Port, pluginPath), "node.plugins", dims)...)

		now := time.Now()
		for _, dp := range dps {
			dp.Timestamp = now
			m.Output.SendDatapoint(dp)
		}

	}, time.Duration(conf.IntervalSeconds)*time.Second)

	return nil
}

func (m *Monitor) fetchNodeInfo(client *http.Client, endpoint string) map[string]string {
	dims := make(map[string]string)

	if nodeJSON, err := getJSON(client, endpoint); err == nil {
		if nodeID, exists := nodeJSON["id"]; exists {
			dims["node_id"], _ = nodeID.(string)
		}
		if nodeName, exists := nodeJSON["name"]; exists {
			dims["node_name"], _ = nodeName.(string)
		}
	}

	return dims
}

func (m *Monitor) fetchMetrics(client *http.Client, endpoint string, prefix string, dimensions map[string]string) []*datapoint.Datapoint {
	var dps []*datapoint.Datapoint

	if metricsJSON, err := getJSON(client, endpoint); err == nil {
		dps = m.extractDatapoints(prefix, metricsJSON, dimensions)
	}

	return dps
}

func (m *Monitor) fetchPipelineMetrics(client *http.Client, endpoint string, prefix string, dimensions map[string]string) []*datapoint.Datapoint {
	var dps []*datapoint.Datapoint

	if metricsJSON, err := getJSON(client, endpoint); err == nil {
		if pipelines, exists := metricsJSON["pipelines"]; exists {
			if pipelinesObj, isObject := pipelines.(map[string]interface{}); isObject {
				for pipelineName, pipelineObj := range pipelinesObj {
					if pipeline, converted := pipelineObj.(map[string]interface{}); converted {
						dimsClone := make(map[string]string)
						for dimk, dimv := range dimensions {
							dimsClone[dimk] = dimv
						}
						dimsClone["pipeline"] = pipelineName
						dps = append(dps, m.extractDatapoints(prefix, pipeline, dimsClone)...)
					}
				}
			}
		}
	}

	return dps
}

func (m *Monitor) fetchPluginMetrics(client *http.Client, endpoint string, prefix string, dimensions map[string]string) []*datapoint.Datapoint {
	if metricsJSON, err := getJSON(client, endpoint); err == nil {
		if total, exists := metricsJSON["total"]; exists {
			metricName := fmt.Sprintf("%s.%s", prefix, "total")
			metricType, isEnabled := m.metricTypeMap[metricName]
			metricValue, castErr := datapoint.CastMetricValueWithBool(total)

			if isEnabled && castErr == nil {
				return []*datapoint.Datapoint{
					&datapoint.Datapoint{
						Metric:     metricName,
						MetricType: metricType,
						Value:      metricValue,
						Dimensions: dimensions,
					},
				}
			}
		}
	}

	return nil
}

func (m *Monitor) extractDatapoints(metricPath string, metricsJSON map[string]interface{}, dims map[string]string) []*datapoint.Datapoint {
	var dps []*datapoint.Datapoint

	for k, v := range metricsJSON {
		childPath := fmt.Sprintf("%s.%s", metricPath, k)
		if obj, isObject := v.(map[string]interface{}); isObject {
			dps = append(dps, m.extractDatapoints(childPath, obj, dims)...)
		} else if arr, isArray := v.([]interface{}); isArray {
			for _, arrayItem := range arr {
				if obj, isObject = arrayItem.(map[string]interface{}); isObject {
					if name, exists := obj["name"]; exists {
						dimsClone := make(map[string]string)
						for dimk, dimv := range dims {
							dimsClone[dimk] = dimv
						}
						dimName, exists := dimensionKeyMap[k]
						if !exists {
							dimName = k
						}
						dimsClone[dimName], _ = name.(string)
						dps = append(dps, m.extractDatapoints(childPath, obj, dimsClone)...)
					}
				}
			}
		} else if metricType, exists := m.metricTypeMap[childPath]; exists {
			if metricValue, err := datapoint.CastMetricValueWithBool(v); err == nil {
				dps = append(dps, &datapoint.Datapoint{
					Metric:     childPath,
					MetricType: metricType,
					Value:      metricValue,
					Dimensions: dims,
				})
			}
		}
	}

	return dps
}

// Shutdown the monitor
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}

func getJSON(client *http.Client, endpoint string) (map[string]interface{}, error) {
	response, err := client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not connect to %s", endpoint)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Could not connect to %s : %s ", endpoint, http.StatusText(response.StatusCode)))
	}

	body, err := ioutil.ReadAll(response.Body)
	metricsJSON := make(map[string]interface{})
	err = json.Unmarshal(body, &metricsJSON)

	return metricsJSON, nil
}
