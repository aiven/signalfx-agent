<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# cpu

Monitor Type: `cpu` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/cpu))

**Accepts Endpoints**: No

**Multiple Instances Allowed**: **No**

## Overview

This monitor reports cpu metrics.

On Linux hosts, this monitor relies on the `/proc` filesystem.
If the underlying host's `/proc` file system is mounted somewhere other than
/proc please specify the path using the top level configuration `procPath`.

```yaml
procPath: /proc
monitors:
 - type: cpu
```


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: cpu
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


This monitor has no configuration options.
## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - ***`cpu.utilization`*** (*gauge*)<br>    Percent of CPU used on this host. This metric is emitted with a plugin dimension set to "signalfx-metadata".
 - `cpu.utilization_per_core` (*gauge*)<br>    Percent of CPU used on each core. This metric is emitted with the plugin dimension set to "signalfx-metadata"

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



