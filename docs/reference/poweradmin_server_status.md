## poweradmin server status

Show the status of the PowerDNS server

### Synopsis

Show whether the PowerDNS server behind Poweradmin is running, its version,
uptime and optionally its metrics and the reachability of the autoprimary servers.

Requires the server_status_view permission (administrators have it implicitly);
--include-slaves additionally requires supermaster_view. The command exits with
a non-zero status when PowerDNS is not reachable, so it can be used in
monitoring checks.

```
poweradmin server status [flags]
```

### Examples

```
  poweradmin server status
  poweradmin server status --metrics uptime,udp-queries
  poweradmin server status --show-metrics -o json
  poweradmin server status --include-slaves
```

### Options

```
  -h, --help              help for status
      --include-slaves    Also probe the autoprimary servers (slower, requires supermaster_view)
      --metrics strings   Only return these metrics (comma-separated, e.g. uptime,udp-queries)
      --no-header         Suppress table header row
  -o, --output string     Output format. One of: table|json|yaml (default "table")
      --show-metrics      Show all metrics in table output (implied by --metrics)
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin server](poweradmin_server.md)	 - Inspect the PowerDNS server

