## poweradmin zones dnssec disable

Unsign a zone

### Synopsis

Remove DNSSEC from a zone and delete its keys.

Remove the DS records from the parent zone first; otherwise validating
resolvers will reject the zone.

```
poweradmin zones dnssec disable [flags]
```

### Examples

```
  poweradmin zones dnssec disable --name example.com --yes
```

### Options

```
  -h, --help            help for disable
      --id string       Zone ID
      --name string     Zone name (e.g. example.com)
  -o, --output string   Output format. One of: table|json|yaml (default "table")
  -y, --yes             Skip confirmation prompt
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones dnssec](poweradmin_zones_dnssec.md)	 - Manage DNSSEC of a zone

