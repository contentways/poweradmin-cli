## poweradmin zones metadata list

List metadata entries for a DNS zone

### Synopsis

List all metadata entries (e.g. ALLOW-AXFR-FROM) for a Poweradmin zone.

```
poweradmin zones metadata list [flags]
```

### Options

```
  -h, --help            help for list
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

* [poweradmin zones metadata](poweradmin_zones_metadata.md)	 - Manage zone metadata

