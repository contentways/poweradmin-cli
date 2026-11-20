## poweradmin zones metadata

Manage zone metadata

### Synopsis

Manage metadata entries (e.g. ALLOW-AXFR-FROM) of a Poweradmin zone.

Without a subcommand the metadata of the zone is listed, like "zones metadata list".

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
* [poweradmin zones metadata delete](poweradmin_zones_metadata_delete.md)	 - Delete zone metadata for a kind
* [poweradmin zones metadata get](poweradmin_zones_metadata_get.md)	 - Get zone metadata by kind
* [poweradmin zones metadata list](poweradmin_zones_metadata_list.md)	 - List metadata entries for a DNS zone
* [poweradmin zones metadata set](poweradmin_zones_metadata_set.md)	 - Set (replace) zone metadata for a kind

