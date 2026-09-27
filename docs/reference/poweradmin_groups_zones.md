## poweradmin groups zones

Manage zones of a group

### Synopsis

Manage the zones assigned to a Poweradmin group.

Without a subcommand the zones of the group are listed, like "groups zones list".

```
poweradmin groups zones [flags]
```

### Options

```
  -h, --help            help for zones
      --id string       Group ID
      --name string     Group name
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

* [poweradmin groups](poweradmin_groups.md)	 - Manage Poweradmin groups
* [poweradmin groups zones add](poweradmin_groups_zones_add.md)	 - Add a zone to a group
* [poweradmin groups zones list](poweradmin_groups_zones_list.md)	 - List zones of a group
* [poweradmin groups zones remove](poweradmin_groups_zones_remove.md)	 - Remove a zone from a group

