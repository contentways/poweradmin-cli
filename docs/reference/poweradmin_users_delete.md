## poweradmin users delete

Delete a user

### Synopsis

Delete a Poweradmin user by username or numeric ID.

Poweradmin refuses to delete a user who still owns zones. Use --transfer-to
(username or numeric ID) to hand those zones over to another user.

```
poweradmin users delete [flags]
```

### Options

```
      --dry-run              Show what would be deleted without making changes
  -h, --help                 help for delete
      --id string            User ID to identify the user
  -i, --interactive          Interactively select users to delete (feature preview)
      --name string          Username to identify the user
  -o, --output string        Output format. One of: table|json|yaml (default "table")
  -q, --quiet                Suppress output after deletion
      --transfer-to string   Username or ID of the user who receives the deleted user's zones
  -y, --yes                  Skip confirmation prompt
```

### Options inherited from parent commands

```
  -k, --api-key string   Poweradmin API key (overrides config and env)
  -u, --url string       Poweradmin URL (e.g. https://dns.example.com)
  -v, --verbose          Log HTTP requests/responses to stderr for debugging
```

### SEE ALSO

* [poweradmin users](poweradmin_users.md)	 - Manage Poweradmin users

