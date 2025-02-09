<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# expvar

Monitor Type: `expvar` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/expvar))

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Overview

The expvar monitor is a SignalFx Smart Agent monitor that scrapes metrics from variables exposed in JSON
format at an HTTP endpoint by [expvar](https://golang.org/pkg/expvar/). The monitor uses configured paths
to get metric and dimension values from fetched JSON objects.

The Metrics section in this document shows metrics derived from expvar variable
[memstats](https://golang.org/pkg/runtime/). The memstat variable is exposed by default. These memstat
metrics are referred to as standard or default metrics. The configuration examples shown are excerpts
limited to the monitor configuration section of the SignalFx Smart Agent configuration file `agent.yml`.

Below is an example showing the minimal required expvar monitor configuration for exporting
the default metrics from endpoint `http://172.17.0.3:8000/debug/vars`. `/debug/vars` is the default path.
```
monitors:
- type: expvar
  host: 172.17.0.3
  path: /debug/vars
  port: 8000
```
We recommend you include the extra dimension `metric_source` with a meaningful value in order to facilitate
filtering in the SignalFx app. See below.
```
monitors:
- type: expvar
  host: 172.17.0.3
  path: /debug/vars
  port: 8000
  extraDimensions:
    metric_source: expvar
```
Below is an example showing part of a JSON payload containing the exposed variable requestsPerSecond containing
requests per second metric information.
```
{
  ...
  "requestsPerSecond": 919,
  ...
}
```
Suppose that the payload is emanating from endpoint `http://172.17.0.4:6000/appmetrics`. The monitor can be
configured as shown below in order to scrape requestsPerSecond. The metric name is optional. If not provided,
the JSONPath value `requestsPerSecond` snake cased to `requests_per_second` will be used instead.
```
monitors:
- type: expvar
  host: 172.17.0.4
  path: /debug/vars
  port: 6000
  metrics:
    - name: requests.sec
      JSONPath: requestsPerSecond
      type: gauge
  extraDimensions:
    metric_source: expvar-aws
```
The expvar monitor can be configured to extract metric values from complex JSON objects such as the one shown
below. Suppose the `memstats` variable shown below is exposed at endpoint `http://172.17.0.5:5000/debug/vars`
and you want to extract the cumulative `Mallocs` values.
```
{
  ...
  "memstats": {
                ...
                "GCCPUFraction": 0.0000032707490586459204,
                "BySize": [
                  {
                      "Size": 32,
                      "Mallocs": 35387,
                      "Frees": 35021
                  },
                  {
                      "Size": 48,
                      "Mallocs": 35387,
                      "Frees": 63283
                  }
                ]
                "HeapAlloc": 2138088,
                ...
              }
  ...
}
```
The should be configured as shown below so as to fetch the `Mallocs` values. The JSONPath is what maps the
location of values with the JSON object. The path must terminate primitive values or an array containing
primitive values. It cannot terminated on embedded object(s).
```
monitors:
- type: expvar
  host: 172.12.0.5
  path: /debug/vars
  port: 5000
  metrics:
    - JSONPath: memstats.BySize.Mallocs
      type: cumulative
  extraDimensions:
    metric_source: expvar
```
No metric name was provided for this configuration so the metric name defaults to memstats.by_size.mallocs.
Because memstats.BySize is an array of size 2 there are 2 values for memstats.BySize.Mallocs (35387 and 35387).
Two data points are created for metric memstats.by_size.mallocs for the 2 values. Additionaly, a dimension
name memstats.by_size containing array index created for each respective datapoint.

Also, custom dimensions can be added to metric as shown below. The dimension name required if a dimension
value is provided whereas it is optional when a JSONPath is provided instead.
```
monitors:
- type: expvar
  host: 172.12.0.5
  path: /debug/vars
  port: 5000
  metrics:
    - JSONPath: memstats.BySize.Mallocs
      type: cumulative
      - dimensions:
        name: physical_memory
        value: 4GiB
      - name: app_mem
        value: "10 MiB"
  extraDimensions:
    metric_source: expvar
```
A dimension JSONPath is configured as shown below. The monitor gets JSON key at the specified path as the
dimension value. The dimension name is optional if the dimension JSONPath is specified. When not provided
the monitor snake cases the dimension JSONPath and uses that for the dimension name. The dimension JSONPath
must be shorter than the metric JSONPath and have the same root.
```
monitors:
- type: expvar
  host: 172.12.0.5
  path: /debug/vars
  port: 5000
  metrics:
    - JSONPath: memstats.BySize.Mallocs
      type: cumulative
      - dimensions:
        JSONPath: memstats
      - dimensions:
        name: by_size_index
        JSONPath: memstats.BySize
  extraDimensions:
    metric_source: expvar
```
DO NOT configure the monitor for memstats metrics because they are standard metrics provided by default.
memstats was used to provide a realistic example.


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: expvar
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `host` | **yes** | `string` | Host of the expvar endpoint |
| `port` | **yes** | `integer` | Port of the expvar endpoint |
| `useHTTPS` | no | `bool` | If true, the agent will connect to the host using HTTPS instead of plain HTTP. (**default:** `false`) |
| `skipVerify` | no | `bool` | If useHTTPS is true and this option is also true, the host's TLS cert will not be verified. (**default:** `false`) |
| `path` | no | `string` | Path to the expvar endpoint, usually `/debug/vars` (the default). (**default:** `/debug/vars`) |
| `enhancedMetrics` | no | `bool` | If true, sends metrics memstats.alloc, memstats.by_size.size, memstats.by_size.mallocs and memstats.by_size.frees (**default:** `false`) |
| `metrics` | no | `list of objects (see below)` | Metrics configurations |


The **nested** `metrics` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `name` | no | `string` | Metric name |
| `JSONPath` | **yes** | `string` | JSON path of the metric value |
| `type` | **yes** | `string` | SignalFx metric type. Possible values are "gauge" or "cumulative" |
| `dimensions` | no | `list of objects (see below)` | Metric dimensions |


The **nested** `dimensions` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `name` | **yes** | `string` | Dimension name |
| `JSONPath` | no | `string` | JSON path of the dimension value |
| `value` | no | `string` | Dimension value |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - `memstats.alloc` (*gauge*)<br>    Bytes of allocated heap objects. Same as memstats.heap_alloc
 - ***`memstats.buck_hash_sys`*** (*gauge*)<br>    Bytes of memory in profiling bucket hash tables
 - `memstats.by_size.frees` (*counter*)<br>    Cumulative count of heap objects freed in a class. The class is identified by dimension class and it is as described for metric memstats.by_size.size
 - `memstats.by_size.mallocs` (*counter*)<br>    Cumulative count of heap objects allocated in a class. The class is identified by dimension class and it is as described for metric memstats.by_size.size
 - `memstats.by_size.size` (*counter*)<br>    The maximum byte size of a class as identified by dimension class. It is the class interval upper limit. The values of dimension class are numbers between 0 and 60 inclusive. Consecutive classes are of consecutive dimension class values. The lower limit of a class is the upper limit of the consecutive class below. Metrics memstats.by_size.size, memstats.by_size.mallocs and memstats.by_size.frees of the same class are related
 - `memstats.debug_gc` (*gauge*)<br>    memstats.debug_gc is currently unused
 - ***`memstats.enable_gc`*** (*gauge*)<br>    Boolean that indicates that GC is enabled. It is always true, even if GOGC=off
 - ***`memstats.frees`*** (*counter*)<br>    Cumulative count of heap objects freed
 - ***`memstats.gc_sys`*** (*gauge*)<br>    Bytes of memory in garbage collection metadata
 - ***`memstats.gccpu_fraction`*** (*gauge*)<br>    The fraction of this program's available CPU time used by the GC since the program started
 - ***`memstats.heap_alloc`*** (*gauge*)<br>    Bytes of allocated heap objects
 - ***`memstats.heap_idle`*** (*gauge*)<br>    Bytes in idle (unused) spans
 - ***`memstats.heap_inuse`*** (*gauge*)<br>    Bytes in in-use spans
 - ***`memstats.heap_objects`*** (*gauge*)<br>    Number of allocated heap objects
 - ***`memstats.heap_released`*** (*gauge*)<br>    Bytes of physical memory returned to the OS
 - `memstats.heap_sys` (*gauge*)<br>    Bytes of heap memory obtained from the OS
 - ***`memstats.last_gc`*** (*gauge*)<br>    The time the last garbage collection finished, as nanoseconds since 1970 (the UNIX epoch)
 - ***`memstats.lookups`*** (*counter*)<br>    Number of pointer lookups performed by the runtime
 - ***`memstats.m_cache_inuse`*** (*gauge*)<br>    Bytes of allocated mcache structures
 - ***`memstats.m_cache_sys`*** (*gauge*)<br>    Bytes of memory obtained from the OS for mcache structures
 - ***`memstats.m_span_inuse`*** (*gauge*)<br>    Bytes of allocated mspan structures
 - ***`memstats.m_span_sys`*** (*gauge*)<br>    Bytes of memory obtained from the OS for mspan
 - ***`memstats.mallocs`*** (*counter*)<br>    Cumulative count of heap objects allocated
 - ***`memstats.most_recent_gc_pause_end`*** (*gauge*)<br>    Most recent GC pause end time, as nanoseconds since 1970 (the UNIX epoch)
 - ***`memstats.most_recent_gc_pause_ns`*** (*gauge*)<br>    Most recent GC stop-the-world pause time in nanoseconds
 - ***`memstats.next_gc`*** (*gauge*)<br>    Target heap size of the next GC cycle
 - ***`memstats.num_forced_gc`*** (*counter*)<br>    Number of GC cycles that were forced by the application calling the GC function
 - ***`memstats.num_gc`*** (*counter*)<br>    Number of completed GC cycles
 - ***`memstats.other_sys`*** (*gauge*)<br>    Bytes of memory in miscellaneous off-heap runtime allocations
 - ***`memstats.pause_total_ns`*** (*counter*)<br>    Cumulative nanoseconds in GC stop-the-world pauses since the program started
 - ***`memstats.stack_inuse`*** (*gauge*)<br>    Bytes in stack spans
 - ***`memstats.stack_sys`*** (*gauge*)<br>    Bytes of stack memory obtained from the OS
 - ***`memstats.sys`*** (*gauge*)<br>    total bytes of memory obtained from the OS
 - ***`memstats.total_alloc`*** (*counter*)<br>    Cumulative bytes allocated for heap objects

### Non-default metrics (version 4.7.0+)

**The following information applies to the agent version 4.7.0+ that has
`enableBuiltInFiltering: true` set on the top level of the agent config.**

To emit metrics that are not _default_, you can add those metrics in the
generic monitor-level `extraMetrics` config option.  Metrics that are derived
from specific configuration options that do not appear in the above list of
metrics do not need to be added to `extraMetrics`.

To see a list of metrics that will be emitted you can run `agent-status
monitors` after configuring this monitor in a running agent instance.

### Legacy non-default metrics (version < 4.7.0)

**The following information only applies to agent version older than 4.7.0. If
you have a newer agent and have set `enableBuiltInFiltering: true` at the top
level of your agent config, see the section above. See upgrade instructions in
[Old-style whitelist filtering](../legacy-filtering.md#old-style-whitelist-filtering).**

If you have a reference to the `whitelist.json` in your agent's top-level
`metricsToExclude` config option, and you want to emit metrics that are not in
that whitelist, then you need to add an item to the top-level
`metricsToInclude` config option to override that whitelist (see [Inclusion
filtering](../legacy-filtering.md#inclusion-filtering).  Or you can just
copy the whitelist.json, modify it, and reference that in `metricsToExclude`.



