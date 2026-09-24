## poweradmin zones metadata

List metadata entries for a DNS zone

### Synopsis

List all metadata entries (e.g. ALLOW-AXFR-FROM) for a Poweradmin zone.

```
poweradmin zones metadata [flags]
```

### Options

```
  -h, --help            help for metadata
      --id string       Zone ID
      --name string     Zone name (e.g. example.com)
      --no-header       Suppress table header row
  -o, --output string   Output format. One of: table|json|yaml (default "table")
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones](poweradmin_zones.md)	 - Manage DNS zones

