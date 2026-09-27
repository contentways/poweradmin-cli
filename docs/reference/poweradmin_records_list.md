## poweradmin records list

List all records in a zone

### Synopsis

List all DNS records in a zone by name or ID.

```
poweradmin records list [flags]
```

### Options

```
  -h, --help               help for list
      --no-header          Suppress table header row
  -o, --output string      Output format. One of: table|full|json|yaml (default "table")
      --sort string        Sort by field. One of: name|type|ttl|prio
      --type string        Filter by record type (e.g. A, AAAA, MX, TXT)
      --zone-id string     Zone ID
      --zone-name string   Zone name (e.g. example.com)
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin records](poweradmin_records.md)	 - Manage DNS records

