## poweradmin zones dnssec keys delete

Delete a DNSSEC key

```
poweradmin zones dnssec keys delete [flags]
```

### Examples

```
  poweradmin zones dnssec keys delete --name example.com --key-id 3 --yes
```

### Options

```
  -h, --help          help for delete
      --id string     Zone ID
      --key-id int    Key ID (required, see "zones dnssec keys list")
      --name string   Zone name (e.g. example.com)
  -q, --quiet         Suppress output after deletion
  -y, --yes           Skip confirmation prompt
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones dnssec keys](poweradmin_zones_dnssec_keys.md)	 - Manage the DNSSEC keys of a zone

