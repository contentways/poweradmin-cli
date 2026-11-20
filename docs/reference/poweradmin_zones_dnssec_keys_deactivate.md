## poweradmin zones dnssec keys deactivate

Deactivate a DNSSEC key

```
poweradmin zones dnssec keys deactivate [flags]
```

### Examples

```
  poweradmin zones dnssec keys deactivate --name example.com --key-id 3
```

### Options

```
  -h, --help            help for deactivate
      --id string       Zone ID
      --key-id int      Key ID (required, see "zones dnssec keys list")
      --name string     Zone name (e.g. example.com)
  -o, --output string   Output format. One of: table|json|yaml (default "table")
  -q, --quiet           Suppress output
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones dnssec keys](poweradmin_zones_dnssec_keys.md)	 - Manage the DNSSEC keys of a zone

