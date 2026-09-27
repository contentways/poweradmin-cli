## poweradmin zones dnssec enable

Sign a zone

### Synopsis

Sign a zone with DNSSEC and show its DS records.

The zone is only validated by resolvers once the DS records are published in
the parent zone (usually at the registrar).

```
poweradmin zones dnssec enable [flags]
```

### Examples

```
  poweradmin zones dnssec enable --name example.com
```

### Options

```
  -h, --help            help for enable
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

