## poweradmin groups members

Manage members of a group

### Synopsis

Manage the users of a Poweradmin group.

Without a subcommand the members of the group are listed, like "groups members list".

```
poweradmin groups members [flags]
```

### Options

```
  -h, --help            help for members
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
* [poweradmin groups members add](poweradmin_groups_members_add.md)	 - Add a user to a group
* [poweradmin groups members list](poweradmin_groups_members_list.md)	 - List members of a group
* [poweradmin groups members remove](poweradmin_groups_members_remove.md)	 - Remove a user from a group

