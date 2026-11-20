## poweradmin zones dnssec keys add

Add a DNSSEC key to a zone

### Synopsis

Add a DNSSEC key to a zone. Keys are created inactive unless --active is set.

Poweradmin validates the combination of algorithm and bits, e.g. ecdsa256
requires 256 bits and the RSA algorithms 1024 or 2048 bits.

```
poweradmin zones dnssec keys add [flags]
```

### Examples

```
  poweradmin zones dnssec keys add --name example.com --type csk --algorithm ecdsa256 --bits 256 --active
  poweradmin zones dnssec keys add --name example.com --type zsk --algorithm rsasha256 --bits 2048
```

### Options

```
      --active             Create the key active
      --algorithm string   Algorithm, e.g. ecdsa256, ed25519, rsasha256 (required)
      --bits int           Key size in bits, e.g. 256 for ecdsa256 (required)
  -h, --help               help for add
      --id string          Zone ID
      --name string        Zone name (e.g. example.com)
  -o, --output string      Output format. One of: table|json|yaml (default "table")
  -q, --quiet              Suppress output
      --type string        Key type: ksk, zsk or csk (required)
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones dnssec keys](poweradmin_zones_dnssec_keys.md)	 - Manage the DNSSEC keys of a zone

