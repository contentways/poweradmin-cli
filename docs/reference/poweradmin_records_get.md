## poweradmin records get

Get a DNS record by ID

### Synopsis

Get a single DNS record by its ID from a zone.

```
poweradmin records get [flags]
```

### Options

```
  -h, --help               help for get
      --id string          Record ID (required)
  -o, --output string      Output format. One of: table|json|yaml (default "table")
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

