<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/df

Monitor Type: `collectd/df` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/df))

**Accepts Endpoints**: No

**Multiple Instances Allowed**: **No**

## Overview

Tracks free disk space on the host using the collectd [df
plugin](https://collectd.org/wiki/index.php/Plugin:DF).

Note that on Linux a filesystem **must** be mounted in the same filesystem
namespace that the agent is running in for this monitor to be able to
collect statistics about that filesystem.  This is mostly an issue when
running the agent in a container.


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: collectd/df
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `hostFSPath` | no | `string` | Path to the root of the host filesystem.  Useful when running in a container and the host filesystem is mounted in some subdirectory under /. |
| `ignoreSelected` | no | `bool` | If true, the filesystems selected by `fsTypes` and `mountPoints` will be excluded and all others included. (**default:** `true`) |
| `fsTypes` | no | `list of strings` | The filesystem types to include/exclude. (**default:** `[aufs overlay tmpfs proc sysfs nsfs cgroup devpts selinuxfs devtmpfs debugfs mqueue hugetlbfs securityfs pstore binfmt_misc autofs]`) |
| `mountPoints` | no | `list of strings` | The mount paths to include/exclude, is interpreted as a regex if surrounded by `/`.  Note that you need to include the full path as the agent will see it, irrespective of the hostFSPath option. (**default:** `[/^/var/lib/docker/ /^/var/lib/rkt/pods/ /^/net// /^/smb//]`) |
| `reportByDevice` | no | `bool` |  (**default:** `false`) |
| `reportInodes` | no | `bool` |  (**default:** `false`) |
| `valuesPercentage` | no | `bool` | If true percent based metrics will be reported. (**default:** `false`) |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - ***`df_complex.free`*** (*gauge*)<br>    Measures free disk space in bytes on this file system.
 - `df_complex.reserved` (*gauge*)<br>    Measures disk space in bytes reserved for the super-user on this file system.
 - ***`df_complex.used`*** (*gauge*)<br>    Measures used disk space in bytes on this file system.
 - `percent_inodes.free` (*gauge*)<br>    Measures free inodes as a percentage of total inodes in the file system.  Inodes are structures used by file systems to store information about files (other than its content).
 - `percent_inodes.reserved` (*gauge*)<br>    Measures inodes reserved for the super-user as a percentage of total inodes in the file system.  Inodes are structures used by file systems to store information about files (other than its content).
 - `percent_inodes.used` (*gauge*)<br>    Measures used inodes as a percentage of total inodes in the file system.  Inodes are structures used by file systems to store information about files (other than its content).

#### Group inodes
All of the following metrics are part of the `inodes` metric group. All of
the non-default metrics below can be turned on by adding `inodes` to the
monitor config option `extraGroups`:
 - `df_inodes.free` (*gauge*)<br>    Measures free inodes in the file system.  Inodes are structures used by Unix filesystems to store metadata about files.
 - `df_inodes.reserved` (*gauge*)<br>    Measures inodes reserved for the super user in the file system.  Inodes are structures used by Unix filesystems to store metadata about files.
 - `df_inodes.used` (*gauge*)<br>    Measures used inodes in the file system.  Inodes are structures used by Unix filesystems to store metadata about files.

#### Group percentage
All of the following metrics are part of the `percentage` metric group. All of
the non-default metrics below can be turned on by adding `percentage` to the
monitor config option `extraGroups`:
 - `percent_bytes.free` (*gauge*)<br>    Measures free disk space as a percentage of total disk space on this file system.
 - `percent_bytes.reserved` (*gauge*)<br>    Measures disk space reserved for the super-user as a percentage of total disk space of this file system.
 - `percent_bytes.used` (*gauge*)<br>    Measures used disk space as a percentage of total disk space of this file system.

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



