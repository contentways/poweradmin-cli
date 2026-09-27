## poweradmin groups zones remove

Remove a zone from a group

### Synopsis

Disassociate a zone from a Poweradmin group.

```
poweradmin groups zones remove [flags]
```

### Options

```
      --group-id string   Group ID (required)
  -h, --help              help for remove
      --zone-id string    Zone ID to remove (required)
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin groups zones](poweradmin_groups_zones.md)	 - Manage zones of a group

