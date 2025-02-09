<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/etcd

Monitor Type: `collectd/etcd` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/etcd))

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Overview

Monitors an etcd key/value store using the [collectd etcd Python plugin](https://github.com/signalfx/collectd-etcd).

Requires etcd 2.0.8 or later.


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: collectd/etcd
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `pythonBinary` | no | `string` | Path to a python binary that should be used to execute the Python code. If not set, a built-in runtime will be used.  Can include arguments to the binary as well. |
| `host` | **yes** | `string` |  |
| `port` | **yes** | `integer` |  |
| `clusterName` | **yes** | `string` | An arbitrary name of the etcd cluster to make it easier to group together and identify instances. |
| `sslKeyFile` | no | `string` | Client private key if using client certificate authentication. |
| `sslCertificate` | no | `string` | Client public key if using client certificate authentication. |
| `sslCACerts` | no | `string` | Certificate authority or host certificate to trust. |
| `skipSSLValidation` | no | `bool` | If `true`, etcd's SSL certificate will not be verified. Enabling this option results in the `sslCACerts` option being ignored. (**default:** `false`) |
| `enhancedMetrics` | no | `bool` |  (**default:** `false`) |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - ***`counter.etcd.leader.counts.fail`*** (*counter*)<br>    Total number of failed rpc requests to with a follower
 - ***`counter.etcd.leader.counts.success`*** (*counter*)<br>    Total number of successful rpc requests to with a follower
 - ***`counter.etcd.self.recvappendreq.cnt`*** (*counter*)<br>    Total number of append requests received by a member
 - ***`counter.etcd.self.sendappendreq.cnt`*** (*counter*)<br>    Total number of append requests sent by a member
 - ***`counter.etcd.store.compareanddelete.fail`*** (*counter*)<br>    Total number of failed compare-and-delete operations
 - ***`counter.etcd.store.compareanddelete.success`*** (*counter*)<br>    Total number of successful compare-and-delete operations
 - ***`counter.etcd.store.compareandswap.fail`*** (*counter*)<br>    Total number of failed compare-and-swap operations
 - ***`counter.etcd.store.compareandswap.success`*** (*counter*)<br>    Total number of successful compare-and-swap operations
 - ***`counter.etcd.store.create.fail`*** (*counter*)<br>    Total number of failed create operations
 - ***`counter.etcd.store.create.success`*** (*counter*)<br>    Total number of successful create operations
 - ***`counter.etcd.store.delete.fail`*** (*counter*)<br>    Total number of failed delete operations
 - ***`counter.etcd.store.delete.success`*** (*counter*)<br>    Total number of successful delete operations
 - ***`counter.etcd.store.expire.count`*** (*counter*)<br>    Total number of items expired due to TTL
 - ***`counter.etcd.store.gets.fail`*** (*counter*)<br>    Total number of failed get operations
 - ***`counter.etcd.store.gets.success`*** (*counter*)<br>    Total number of successful get operations
 - ***`counter.etcd.store.sets.fail`*** (*counter*)<br>    Total number of failed set operations
 - ***`counter.etcd.store.sets.success`*** (*counter*)<br>    Total number of successful set operations
 - ***`counter.etcd.store.update.fail`*** (*counter*)<br>    Total number of failed update operations
 - ***`counter.etcd.store.update.success`*** (*counter*)<br>    Total number of successful update operations
 - `gauge.etcd.leader.latency.average` (*gauge*)<br>    Average latency of a follower with respect to the leader
 - ***`gauge.etcd.leader.latency.current`*** (*gauge*)<br>    Current latency of a follower with respect to the leader
 - `gauge.etcd.leader.latency.max` (*gauge*)<br>    Max latency of a follower with respect to the leader
 - `gauge.etcd.leader.latency.min` (*gauge*)<br>    Min latency of a follower with respect to the leader
 - `gauge.etcd.leader.latency.stddev` (*gauge*)<br>    Std dev latency of a follower with respect to the leader
 - ***`gauge.etcd.self.recvbandwidth.rate`*** (*gauge*)<br>    Bandwidth rate of a follower
 - ***`gauge.etcd.self.recvpkg.rate`*** (*gauge*)<br>    Rate at which a follower receives packages
 - ***`gauge.etcd.self.sendbandwidth.rate`*** (*gauge*)<br>    Bandwidth rate of a leader
 - ***`gauge.etcd.self.sendpkg.rate`*** (*gauge*)<br>    Rate at which a leader sends packages
 - ***`gauge.etcd.store.watchers`*** (*gauge*)<br>    Number of watchers

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



