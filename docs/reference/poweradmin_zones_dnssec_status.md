## poweradmin zones dnssec status

Show the DNSSEC status and DS records of a zone

```
poweradmin zones dnssec status [flags]
```

### Examples

```
  poweradmin zones dnssec status --name example.com
  poweradmin zones dnssec status --name example.com -o json
```

### Options

```
  -h, --help            help for status
      --id string       Zone ID
      --name string     Zone name (e.g. example.com)
  -o, --output string   Output format. One of: table|json|yaml (default "table")
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones dnssec](poweradmin_zones_dnssec.md)	 - Manage DNSSEC of a zone

