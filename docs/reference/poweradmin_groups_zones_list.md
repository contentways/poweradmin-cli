## poweradmin groups zones list

List zones of a group

### Synopsis

List all zones associated with a Poweradmin group.

```
poweradmin groups zones list [flags]
```

### Options

```
  -h, --help            help for list
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

* [poweradmin groups zones](poweradmin_groups_zones.md)	 - Manage zones of a group

