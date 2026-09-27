## poweradmin zones dnssec

Manage DNSSEC of a zone

### Synopsis

Sign and unsign zones, show their DS records, manage their DNSSEC keys and
rectify them.

Key management and rectify require Poweradmin 4.5 or newer. Changing DNSSEC
requires the zone_dnssec_manage_own permission for the zone.

### Options

```
  -h, --help   help for dnssec
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin zones](poweradmin_zones.md)	 - Manage DNS zones
* [poweradmin zones dnssec disable](poweradmin_zones_dnssec_disable.md)	 - Unsign a zone
* [poweradmin zones dnssec enable](poweradmin_zones_dnssec_enable.md)	 - Sign a zone
* [poweradmin zones dnssec keys](poweradmin_zones_dnssec_keys.md)	 - Manage the DNSSEC keys of a zone
* [poweradmin zones dnssec rectify](poweradmin_zones_dnssec_rectify.md)	 - Rectify a signed zone
* [poweradmin zones dnssec status](poweradmin_zones_dnssec_status.md)	 - Show the DNSSEC status and DS records of a zone

