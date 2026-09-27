## poweradmin zones dnssec rectify

Rectify a signed zone

### Synopsis

Recalculate the DNSSEC ordering and auth fields of a signed zone
(like pdnsutil rectify-zone). Not possible for unsigned, presigned or
secondary zones. Requires Poweradmin 4.5 or newer.

```
poweradmin zones dnssec rectify [flags]
```

### Examples

```
  poweradmin zones dnssec rectify --name example.com
```

### Options

```
  -h, --help          help for rectify
      --id string     Zone ID
      --name string   Zone name (e.g. example.com)
  -q, --quiet         Suppress output
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones dnssec](poweradmin_zones_dnssec.md)	 - Manage DNSSEC of a zone

