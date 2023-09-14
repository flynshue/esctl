# esctl<!-- omit in toc -->
CLI tool written in Go to for interacting with elasticsearch API. 

This tool aims to be the GO replacement of https://github.com/slmingol/escli.

:memo: **Note:** not all commands have been migrated over yet.

- [Configuration](#configuration)
- [Usage](#usage)
  - [Show cluster health](#show-cluster-health)
  - [List Shard count for each node](#list-shard-count-for-each-node)
  - [Show shard routing allocations](#show-shard-routing-allocations)
  - [Disable shard routing allocations](#disable-shard-routing-allocations)
  - [Enable shard routing allocations](#enable-shard-routing-allocations)
  - [List nodes and their disk usage](#list-nodes-and-their-disk-usage)
  - [List nodes and usage statistics](#list-nodes-and-usage-statistics)
  - [List nodes and their elastic search version](#list-nodes-and-their-elastic-search-version)
  - [List all indexes and their sizes sorted (big -\> small)](#list-all-indexes-and-their-sizes-sorted-big---small)
  - [List all indexes and their version](#list-all-indexes-and-their-version)
  - [List index templates with their index patterns](#list-index-templates-with-their-index-patterns)
  - [Get more details for index template](#get-more-details-for-index-template)
  - [List all shards in cluster and their node](#list-all-shards-in-cluster-and-their-node)
  - [List all shards on specific node](#list-all-shards-on-specific-node)
  - [Watch active shard recovery](#watch-active-shard-recovery)

## Configuration
By default, the tool will use `$USER` as the elasticsearch username.  To override that:

```bash
export ES_USERNAME=fooBar
```

To set the elasticsearch password:
```bash
export ES_PASSWORD=$(lpass show --password example.com)
```

To set the hosts:
```bash
export ES_HOSTS="https://es-data-01d.lab1.example:9200 https://es-data-01e.lab1.example:9200"
```

You can also use a configuration file instead.
```bash
vim ~/.esctl-lab.yaml 
```

```yaml
# ~/.esctl-lab.yaml 
hosts:
  - https://es-data-01d.lab1.example:9200
  - https://es-data-01e.lab1.example:9200
  - https://es-data-01f.lab1.example:9200
  - https://es-data-01g.lab1.example:9200
# skip tls verification, default is false
insecure: true
# not recommended but you can do it
username: fooBar
password: fakePass
```

## Usage
Every command has built in help and should give a pretty good idea on how to to use it.

Here's some example usage

### Show cluster health
```bash
$ ./esctl --config ~/.esctl-dev.yaml get health
{
  "cluster_name" : "docker-cluster",
  "status" : "green",
  "timed_out" : false,
  "number_of_nodes" : 3,
  "number_of_data_nodes" : 3,
  "active_primary_shards" : 22,
  "active_shards" : 45,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 0,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 100.0
}
```

### List Shard count for each node
```bash
$ ./esctl --config ~/.esctl-dev.yaml list shard count
name         shard_count  
es-data-02   15           
es-data-03   15           
es-data-01   15
```

### Show shard routing allocations
```bash
$ ./esctl --config ~/.esctl-dev.yaml get shards allocations 
defaults
cluster.routing.allocation.enable: all
```

### Disable shard routing allocations
```bash
$ ./esctl --config ~/.esctl-dev.yaml disable shard alloc
{"acknowledged":true,"persistent":{},"transient":{"cluster":{"routing":{"allocation":{"enable":"none"}}}}}
```

### Enable shard routing allocations
```bash
$ ./esctl --config ~/.esctl-dev.yaml enable shard alloc
{"acknowledged":true,"persistent":{},"transient":{"cluster":{"routing":{"allocation":{"enable":"all"}}}}}
```

### List nodes and their disk usage
```bash
$ ./esctl --config ~/.esctl-dev.yaml list nodes storage
ip         node.role   master name       disk.total disk.used disk.avail disk.used_percent
172.18.0.4 cdfhilmrstw *      es-data-02    464.3gb    87.5gb    376.8gb             18.84
172.18.0.2 cdfhilmrstw -      es-data-01    464.3gb    87.5gb    376.8gb             18.84
172.18.0.5 cdfhilmrstw -      es-data-03    464.3gb    87.5gb    376.8gb             18.84

valid data node suffixes: 03, 01, 02
total data nodes: 3
```

### List nodes and usage statistics
```bash
$ ./esctl --config ~/.esctl-dev.yaml list nodes stats
ip         name       heap.percent ram.percent cpu load_1m load_5m load_15m node.role   master disk.total disk.used disk.avail disk.used_percent
172.18.0.2 es-data-01           38         100  10    1.20    1.01     1.08 cdfhilmrstw -         464.3gb    87.5gb    376.8gb             18.84
172.18.0.4 es-data-02           36         100  10    1.20    1.01     1.08 cdfhilmrstw *         464.3gb    87.5gb    376.8gb             18.84
172.18.0.5 es-data-03           70         100  10    1.20    1.01     1.08 cdfhilmrstw -         464.3gb    87.5gb    376.8gb             18.84

valid data node suffixes: 03, 02, 01
```

### List nodes and their elastic search version
```bash
$ ./esctl --config ~/.esctl-dev.yaml list nodes versions
node         elastic-version   ip      roles                                                                                              
es-data-01   172.18.0.2        8.9.1   datadata_colddata_contentdata_frozendata_hotdata_warmingestmastermlremote_cluster_clienttransform  
es-data-02   172.18.0.4        8.9.1   datadata_colddata_contentdata_frozendata_hotdata_warmingestmastermlremote_cluster_clienttransform  
es-data-03   172.18.0.5        8.9.1   datadata_colddata_contentdata_frozendata_hotdata_warmingestmastermlremote_cluster_clienttransform  
flynshue@flynshue-Latitude-7430:~/github.com/flynshue/esctl$ 
```

### List all indexes and their sizes sorted (big -> small)
```bash
$ ./esctl --config ~/.esctl-dev.yaml list index sizes
index                                                        pri rep docs.count store.size pri.store.size
.internal.alerts-observability.uptime.alerts-default-000001    1   1          0          0              0
.fleet-files-agent-000001                                      1   1          0          0              0
.internal.alerts-observability.slo.alerts-default-000001       1   1          0          0              0
.internal.alerts-security.alerts-default-000001                1   1          0          0              0
.fleet-file-data-agent-000001                                  1   1          0          0              0
.internal.alerts-observability.metrics.alerts-default-000001   1   1          0          0              0
.internal.alerts-observability.logs.alerts-default-000001      1   1          0          0              0
.internal.alerts-observability.apm.alerts-default-000001       1   1          0          0              0
```

### List all indexes and their version
```bash
$ ./esctl --config ~/.esctl-dev.yaml list index versions
index                                                           version  
.monitoring-logstash-7-2023.08.18                               7.17.5    
.reporting-2020.09.06                                           7.6.2    
.reporting-2020.10.04                                           7.6.2 
```

### List index templates with their index patterns
```bash
$ ./esctl --config ~/.esctl-dev.yaml list index templates
name                                                        index_patterns                                            order      version composed_of
.alerts-observability.apm.alerts-default-index-template     [.internal.alerts-observability.apm.alerts-default-*]     7                  [.alerts-observability.apm.alerts-mappings, .alerts-legacy-alert-mappings, .alerts-framework-mappings]
.alerts-observability.logs.alerts-default-index-template    [.internal.alerts-observability.logs.alerts-default-*]    7                  [.alerts-ecs-mappings, .alerts-observability.logs.alerts-mappings, .alerts-legacy-alert-mappings, .alerts-framework-mappings]
.alerts-observability.metrics.alerts-default-index-template [.internal.alerts-observability.metrics.alerts-default-*] 7                  [.alerts-ecs-mappings, .alerts-observability.metrics.alerts-mappings, .alerts-legacy-alert-mappings, .alerts-framework-mappings]
.alerts-observability.slo.alerts-default-index-template     [.internal.alerts-observability.slo.alerts-default-*]     7                  [.alerts-observability.slo.alerts-mappings, .alerts-legacy-alert-mappings, .alerts-framework-mappings]
.alerts-observability.uptime.alerts-default-index-template  [.internal.alerts-observability.uptime.alerts-default-*]  7                  [.alerts-observability.uptime.alerts-mappings, .alerts-legacy-alert-mappings, .alerts-framework-mappings]
.alerts-security.alerts-default-index-template              [.internal.alerts-security.alerts-default-*]              7                  [.alerts-ecs-mappings, .alerts-security.alerts-mappings, .alerts-legacy-alert-mappings, .alerts-framework-mappings]
.deprecation-indexing-template                              [.logs-deprecation.*]                                     1000       1       [.deprecation-indexing-mappings, .deprecation-indexing-settings]
.fleet-file-data                                            [.fleet-file-data-*-*]                                    200        1       []
.fleet-filedelivery-data                                    [.fleet-filedelivery-data-*-*]                            200        1       []
.fleet-filedelivery-meta                                    [.fleet-filedelivery-meta-*-*]                            200        1       []
.fleet-files                                                [.fleet-files-*-*]                                        200        1       []
.kibana-event-log-8.9.1-template                            [.kibana-event-log-8.9.1]                                 50                 []
.ml-anomalies-                                              [.ml-anomalies-*]                                         2147483647 8090199 []
.ml-notifications-000002                                    [.ml-notifications-000002]                                2147483647 8090199 []
.ml-state                                                   [.ml-state*]                                              2147483647 8090199 []
.ml-stats                                                   [.ml-stats-*]                                             2147483647 8090199 []
.monitoring-alerts-7                                        [.monitoring-alerts-7]                                    0          8080099 
.monitoring-beats                                           [.monitoring-beats-7-*]                                   0          8080099 
.monitoring-beats-mb                                        [.monitoring-beats-8-*]                                   0          8000108 []
.monitoring-ent-search-mb                                   [.monitoring-ent-search-8-*]                              0          8000108 []
.monitoring-es                                              [.monitoring-es-7-*]                                      0          8080099 
.monitoring-es-mb                                           [.monitoring-es-8-*]                                      0          8000108 []
.monitoring-kibana                                          [.monitoring-kibana-7-*]                                  0          8080099 
.monitoring-kibana-mb                                       [.monitoring-kibana-8-*]                                  0          8000108 []
.monitoring-logstash                                        [.monitoring-logstash-7-*]                                0          8080099 
.monitoring-logstash-mb                                     [.monitoring-logstash-8-*]                                0          8000108 []
.slm-history                                                [.slm-history-5*]                                         2147483647 5       []
.watch-history-16                                           [.watcher-history-16*]                                    2147483647 16      []
apm-source-map                                              [.apm-source-map]                                         0          1       []
behavioral_analytics-events-default                         [behavioral_analytics-events-*]                           100        2       [behavioral_analytics-events-settings, behavioral_analytics-events-mappings]
ilm-history                                                 [ilm-history-5*]                                          2147483647 5       []
logs                                                        [logs-*-*]                                                100        3       [logs-mappings, logs-settings, logs@custom, ecs@dynamic_templates]
metrics                                                     [metrics-*-*]                                             100        3       [metrics-mappings, data-streams-mappings, metrics-settings]
synthetics                                                  [synthetics-*-*]                                          100        3       [synthetics-mappings, data-streams-mappings, synthetics-settings]
synthetics-browser                                          [synthetics-browser-*]                                    200                [synthetics-browser@package, synthetics-browser@custom, .fleet_globals-1, .fleet_agent_id_verification-1]
synthetics-browser.network                                  [synthetics-browser.network-*]                            200                [synthetics-browser.network@package, synthetics-browser.network@custom, .fleet_globals-1, .fleet_agent_id_verification-1]
synthetics-browser.screenshot                               [synthetics-browser.screenshot-*]                         200                [synthetics-browser.screenshot@package, synthetics-browser.screenshot@custom, .fleet_globals-1, .fleet_agent_id_verification-1]
synthetics-http                                             [synthetics-http-*]                                       200                [synthetics-http@package, synthetics-http@custom, .fleet_globals-1, .fleet_agent_id_verification-1]
synthetics-icmp                                             [synthetics-icmp-*]                                       200                [synthetics-icmp@package, synthetics-icmp@custom, .fleet_globals-1, .fleet_agent_id_verification-1]
synthetics-tcp                                              [synthetics-tcp-*]                                        200                [synthetics-tcp@package, synthetics-tcp@custom, .fleet_globals-1, .fleet_agent_id_verification-1]
```

**Index templates (legacy) are deprecated and will be replaced with composable templates**

Here's a quick way to find the legacy index templates on your cluster
```bash
$ ./esctl --config ~/.esctl-dev.yaml list index templates --legacy
name                   index_pattern                order  
.monitoring-logstash   [.monitoring-logstash-7-*]   0      
.monitoring-kibana     [.monitoring-kibana-7-*]     0      
.monitoring-alerts-7   [.monitoring-alerts-7]       0      
.monitoring-es         [.monitoring-es-7-*]         0      
.monitoring-beats      [.monitoring-beats-7-*]      0  
```

### Get more details for index template
```bash
l$ ./esctl --config ~/.esctl-dev.yaml get index template logs
{
  "index_templates" : [
    {
      "name" : "logs",
      "index_template" : {
        "index_patterns" : [
          "logs-*-*"
        ],
        "composed_of" : [
          "logs-mappings",
          "logs-settings",
          "logs@custom",
          "ecs@dynamic_templates"
        ],
        "priority" : 100,
        "version" : 3,
        "_meta" : {
          "description" : "default logs template installed by x-pack",
          "managed" : true
        },
        "data_stream" : {
          "hidden" : false,
          "allow_custom_routing" : false
        },
        "allow_auto_create" : true,
        "ignore_missing_component_templates" : [
          "logs@custom"
        ]
      }
    }
  ]
}
```

More details for a legacy index template
```bash
$ ./esctl --config ~/.esctl-dev.yaml get index template .monitoring-logstash
Warning: .monitoring-logstash is a legacy index template. Legacy index templates have been deprecated starting in 7.8
{
  ".monitoring-logstash" : {
    "order" : 0,
    "version" : 8080099,
    "index_patterns" : [
      ".monitoring-logstash-7-*"
    ],
    "settings" : {
      "index" : {
        "format" : "7",
        "codec" : "best_compression",
        "number_of_shards" : "1",
        "auto_expand_replicas" : "0-1",
        "number_of_replicas" : "0"
      }
    },
    "mappings" : {
      "dynamic" : false,
      "properties" : {
        "cluster_uuid" : {
          "type" : "keyword"
```
The warning messages are written to stderr so you that you can still output the settings to json like below
```bash
$ ./esctl --config ~/.esctl-dev.yaml get index template .monitoring-logstash > /tmp/monitoring-logstash.json
Warning: .monitoring-logstash is a legacy index template. Legacy index templates have been deprecated starting in 7.8
```

```bash
$ cat /tmp/monitoring-logstash.json | jq '.[]|keys'
[
  "aliases",
  "index_patterns",
  "mappings",
  "order",
  "settings",
  "version"
]
```

### List all shards in cluster and their node
```bash
$ ./esctl --config ~/.esctl-dev.yaml list shards
index                                                         shard prirep state   docs   store ip         node
.kibana_analytics_8.9.1_001                                   0     p      STARTED    5   2.3mb 172.19.0.3 es-data-01
.kibana_analytics_8.9.1_001                                   0     r      STARTED    5   2.3mb 172.19.0.5 es-data-03
.security-7                                                   0     r      STARTED  134 462.2kb 172.19.0.3 es-data-01
.security-7                                                   0     p      STARTED  134 419.8kb 172.19.0.5 es-data-03
.kibana_ingest_8.9.1_001                                      0     p      STARTED  136 283.3kb 172.19.0.4 es-data-02
.kibana_task_manager_8.9.1_001                                0     r      STARTED   25 244.2kb 172.19.0.5 es-data-03
.kibana_ingest_8.9.1_001                                      0     r      STARTED   35 244.2kb 172.19.0.3 es-data-01
.kibana_task_manager_8.9.1_001                                0     p      STARTED   25   238kb 172.19.0.4 es-data-02
.kibana_8.9.1_001                                             0     r      STARTED   12  56.9kb 172.19.0.4 es-data-02
.kibana_8.9.1_001                                             0     p      STARTED   12  56.9kb 172.19.0.3 es-data-01
.ds-ilm-history-5-2023.09.14-000001                           0     r      STARTED   27  23.2kb 172.19.0.3 es-data-01
.ds-ilm-history-5-2023.09.14-000001                           0     p      STARTED   27  23.1kb 172.19.0.5 es-data-03
.ds-.logs-deprecation.elasticsearch-default-2023.09.14-000001 0     r      STARTED    1  10.9kb 172.19.0.4 es-data-02
.ds-.logs-deprecation.elasticsearch-default-2023.09.14-000001 0     p      STARTED    1  10.8kb 172.19.0.3 es-data-01
.kibana_alerting_cases_8.9.1_001                              0     p      STARTED    1   6.6kb 172.19.0.5 es-data-03
.kibana_alerting_cases_8.9.1_001                              0     r      STARTED    1   6.6kb 172.19.0.4 es-data-02
.ds-.kibana-event-log-8.9.1-2023.09.14-000001                 0     r      STARTED    1   6.2kb 172.19.0.3 es-data-01
.ds-.kibana-event-log-8.9.1-2023.09.14-000001                 0     p      STARTED    1   6.1kb 172.19.0.5 es-data-03
.apm-agent-configuration                                      0     p      STARTED    0    225b 172.19.0.5 es-data-03
.apm-agent-configuration                                      0     r      STARTED    0    225b 172.19.0.3 es-data-01
.apm-custom-link                                              0     p      STARTED    0    225b 172.19.0.3 es-data-01
.apm-custom-link                                              0     r      STARTED    0    225b 172.19.0.4 es-data-02
.apm-source-map                                               0     r      STARTED    0    225b 172.19.0.5 es-data-03
.apm-source-map                                               0     p      STARTED    0    225b 172.19.0.3 es-data-01
.apm-source-map                                               0     r      STARTED    0    225b 172.19.0.4 es-data-02
.fleet-file-data-agent-000001                                 0     r      STARTED    0    225b 172.19.0.5 es-data-03
.fleet-file-data-agent-000001                                 0     p      STARTED    0    225b 172.19.0.4 es-data-02
.fleet-files-agent-000001                                     0     p      STARTED    0    225b 172.19.0.3 es-data-01
.fleet-files-agent-000001                                     0     r      STARTED    0    225b 172.19.0.4 es-data-02
.internal.alerts-observability.apm.alerts-default-000001      0     p      STARTED    0    225b 172.19.0.3 es-data-01
.internal.alerts-observability.apm.alerts-default-000001      0     r      STARTED    0    225b 172.19.0.4 es-data-02
.internal.alerts-observability.logs.alerts-default-000001     0     p      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.logs.alerts-default-000001     0     r      STARTED    0    225b 172.19.0.4 es-data-02
.internal.alerts-observability.metrics.alerts-default-000001  0     p      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.metrics.alerts-default-000001  0     r      STARTED    0    225b 172.19.0.3 es-data-01
.internal.alerts-observability.slo.alerts-default-000001      0     p      STARTED    0    225b 172.19.0.3 es-data-01
.internal.alerts-observability.slo.alerts-default-000001      0     r      STARTED    0    225b 172.19.0.4 es-data-02
.internal.alerts-observability.uptime.alerts-default-000001   0     r      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.uptime.alerts-default-000001   0     p      STARTED    0    225b 172.19.0.4 es-data-02
.internal.alerts-security.alerts-default-000001               0     r      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-security.alerts-default-000001               0     p      STARTED    0    225b 172.19.0.4 es-data-02
.kibana_security_session_1                                    0     r      STARTED    0    225b 172.19.0.5 es-data-03
.kibana_security_session_1                                    0     p      STARTED    0    225b 172.19.0.4 es-data-02
.kibana_security_solution_8.9.1_001                           0     p      STARTED    0    225b 172.19.0.5 es-data-03
.kibana_security_solution_8.9.1_001                           0     r      STARTED    0    225b 172.19.0.3 es-data-01
```

### List all shards on specific node
```bash
$ ./esctl --config ~/.esctl-dev.yaml list shards --node es-data-03
.kibana_analytics_8.9.1_001                                   0     r      STARTED    5   2.3mb 172.19.0.5 es-data-03
.security-7                                                   0     p      STARTED  134 419.8kb 172.19.0.5 es-data-03
.kibana_task_manager_8.9.1_001                                0     r      STARTED   25 181.9kb 172.19.0.5 es-data-03
.ds-ilm-history-5-2023.09.14-000001                           0     p      STARTED   27  23.1kb 172.19.0.5 es-data-03
.kibana_alerting_cases_8.9.1_001                              0     p      STARTED    1   6.6kb 172.19.0.5 es-data-03
.ds-.kibana-event-log-8.9.1-2023.09.14-000001                 0     p      STARTED    1   6.1kb 172.19.0.5 es-data-03
.apm-agent-configuration                                      0     p      STARTED    0    225b 172.19.0.5 es-data-03
.apm-source-map                                               0     r      STARTED    0    225b 172.19.0.5 es-data-03
.fleet-file-data-agent-000001                                 0     r      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.logs.alerts-default-000001     0     p      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.metrics.alerts-default-000001  0     p      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.uptime.alerts-default-000001   0     r      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-security.alerts-default-000001               0     r      STARTED    0    225b 172.19.0.5 es-data-03
.kibana_security_session_1                                    0     r      STARTED    0    225b 172.19.0.5 es-data-03
.kibana_security_solution_8.9.1_001                           0     p      STARTED    0    225b 172.19.0.5 es-data-03
```

You can sort by shard size too, the default is descending
```bash
$ ./esctl --config ~/.esctl-dev.yaml list shards --node es-data-03 -s asc
.apm-agent-configuration                                      0     p      STARTED    0    225b 172.19.0.5 es-data-03
.apm-source-map                                               0     r      STARTED    0    225b 172.19.0.5 es-data-03
.fleet-file-data-agent-000001                                 0     r      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.logs.alerts-default-000001     0     p      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.metrics.alerts-default-000001  0     p      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-observability.uptime.alerts-default-000001   0     r      STARTED    0    225b 172.19.0.5 es-data-03
.internal.alerts-security.alerts-default-000001               0     r      STARTED    0    225b 172.19.0.5 es-data-03
.kibana_security_session_1                                    0     r      STARTED    0    225b 172.19.0.5 es-data-03
.kibana_security_solution_8.9.1_001                           0     p      STARTED    0    225b 172.19.0.5 es-data-03
.ds-.kibana-event-log-8.9.1-2023.09.14-000001                 0     p      STARTED    1   6.1kb 172.19.0.5 es-data-03
.kibana_alerting_cases_8.9.1_001                              0     p      STARTED    1   6.6kb 172.19.0.5 es-data-03
.ds-ilm-history-5-2023.09.14-000001                           0     p      STARTED   27  23.1kb 172.19.0.5 es-data-03
.kibana_task_manager_8.9.1_001                                0     r      STARTED   25 139.2kb 172.19.0.5 es-data-03
.security-7                                                   0     p      STARTED  134 419.8kb 172.19.0.5 es-data-03
.kibana_analytics_8.9.1_001                                   0     r      STARTED    5   2.3mb 172.19.0.5 es-data-03
```

### Watch active shard recovery
```bash
$ ./esctl --config ~/.esctl-dev.yaml top recov
index                                                    shard time  type stage source_node target_node files files_recovered files_percent bytes_total bytes_percent translog_ops_recovered translog_ops translog_ops_percent
.fleet-files-agent-000001                                0     197ms peer init  es-data-01  es-data-02  0     0               0.0%          0           0.0%          0                      -1           -1.0%
.internal.alerts-observability.apm.alerts-default-000001 0     160ms peer index es-data-02  es-data-01  1     1               100.0%        0           100.0%        0                      0            100.0%



Hit enter to stop
```