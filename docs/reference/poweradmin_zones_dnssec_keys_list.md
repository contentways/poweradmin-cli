## poweradmin zones dnssec keys list

List the DNSSEC keys of a zone

```
poweradmin zones dnssec keys list [flags]
```

### Examples

```
  poweradmin zones dnssec keys list --name example.com
  poweradmin zones dnssec keys list --name example.com -o json
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

* [poweradmin zones dnssec keys](poweradmin_zones_dnssec_keys.md)	 - Manage the DNSSEC keys of a zone

