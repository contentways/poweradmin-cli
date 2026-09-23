# poweradmin-cli

[![CI](https://github.com/contentways/poweradmin-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/contentways/poweradmin-cli/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/contentways/poweradmin-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/contentways/poweradmin-cli)

A command-line interface for managing DNS zones, records, users and groups via the [Poweradmin](https://www.poweradmin.org) REST API.

Built with [poweradmin-go](https://github.com/contentways/poweradmin-go) — the Go SDK for Poweradmin.

## Requirements

- Poweradmin 4.3.0+ running in API-mode (`PA_DNS_BACKEND=api`)
- A valid Poweradmin API key

## Installation

### From release

Download the latest binary for your platform from the [releases page](https://github.com/contentways/poweradmin-cli/releases):

```bash
# Linux (amd64)
curl -L https://github.com/contentways/poweradmin-cli/releases/latest/download/poweradmin-cli_Linux_x86_64.tar.gz | tar xz
sudo mv poweradmin /usr/local/bin/
```

### Via go install

```bash
go install github.com/contentways/poweradmin-cli/v3@latest
```

### Docker

```bash
docker pull contentwaysorg/poweradmin-cli:latest

docker run --rm \
  -e POWERADMIN_URL=https://dns.example.com \
  -e POWERADMIN_API_KEY=pwa_... \
  contentwaysorg/poweradmin-cli:latest zones list
```

### From source

```bash
git clone git@github.com:contentways/poweradmin-cli.git
cd poweradmin-cli
make build
```

The binary is placed at `build/poweradmin`.

## Configuration

Credentials are resolved in the following order of precedence:

| Source | Example |
|--------|---------|
| CLI flags | `-u`, `-k` |
| Environment variables | `POWERADMIN_URL`, `POWERADMIN_API_KEY` |
| Config file | `~/.config/poweradmin/config.yaml` |

### Config file

```yaml
url: https://dns.example.com
api_key: pwa_your_api_key_here
```

Create the config directory and file:

```bash
mkdir -p ~/.config/poweradmin
cat > ~/.config/poweradmin/config.yaml << 'EOF'
url: https://dns.example.com
api_key: pwa_your_api_key_here
EOF
```

### Environment variables

```bash
export POWERADMIN_URL=https://dns.example.com
export POWERADMIN_API_KEY=pwa_your_api_key_here
```

### CLI flags

```bash
poweradmin zones list -u https://dns.example.com -k pwa_...
```

## Shell Completion

Shell completion supports Tab completion for command names, flags and values
including zone names, usernames and group names fetched live from the API.

### zsh (Oh My Zsh)

```bash
poweradmin completion zsh > ~/.oh-my-zsh/completions/_poweradmin
source ~/.zshrc
```

### bash

```bash
poweradmin completion bash > ~/.bash_completion.d/poweradmin
source ~/.bash_completion.d/poweradmin
```

## Usage

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--url` | `-u` | Poweradmin URL |
| `--api-key` | `-k` | Poweradmin API key |
| `--output` | `-o` | Output format: `table`, `full`, `json`, `yaml` |
| `--no-header` | | Suppress table header row |
| `--quiet` | `-q` | Only print ID (create) or suppress output (delete) |
| `--yes` | `-y` | Skip delete confirmation prompt |
| `--verbose` | `-v` | Log HTTP requests/responses to stderr for debugging |
| `--dry-run` | | Show what would be deleted without making changes (delete commands) |

### Verbose Mode

The `--verbose` (`-v`) flag logs each HTTP request to stderr — method,
path, status code, and duration — useful for debugging connection issues
or unexpected API responses without affecting stdout, so it's safe to
combine with scripting and piping.

```bash
poweradmin --verbose zones list
```

```
[poweradmin] GET zones?page=1&per_page=100 -> 200 (126ms)
ID    NAME                 TYPE
...
```

### Interactive Mode

> **Feature preview:** Interactive mode is new and its prompts/flow may
> change in a future release. Feedback is welcome.

The `create` commands for zones, users, groups, permission templates and
records support an `--interactive` (`-i`) flag. When set, any required
value not already supplied via flag or argument is prompted for
interactively, followed by a summary and confirmation before the resource
is created.

```bash
poweradmin zones create --interactive
```

The `delete` commands support `--interactive` as well, but with a
different flow: instead of prompting for identifying flags, all
resources are listed and you pick zero or more via a multi-select
(arrow keys to navigate, space to toggle, enter to confirm), followed
by a summary and confirmation before deletion.

```bash
poweradmin zones delete --interactive
```

### Dry Run Mode

All `delete` commands support a `--dry-run` flag. It lists what would be
deleted and exits without prompting for confirmation or making any API
calls — useful for previewing the effect of a name/ID match or an
interactive multi-select before committing to it.

```bash
poweradmin zones delete --name example.com --dry-run
poweradmin zones delete --interactive --dry-run
```

### Zones

```bash
# List all zones
poweradmin zones list
poweradmin zones list -o json

# Filter and sort zones
poweradmin zones list --type NATIVE
poweradmin zones list --name-filter contentways
poweradmin zones list --sort name

# Get a zone by name or ID
poweradmin zones get --name example.com
poweradmin zones get --id 42 -o json

# Create a zone with nameservers
poweradmin zones create example.com
poweradmin zones create example.com \
  --nameserver ns1.example.com,ns2.example.com,ns3.example.com

# Create a zone interactively
poweradmin zones create --interactive

# Create and capture the ID
ZONE_ID=$(poweradmin zones create example.com -q)

# Delete a zone
poweradmin zones delete --name example.com
poweradmin zones delete --name example.com --yes

# Delete zones interactively (multi-select)
poweradmin zones delete --interactive

# Preview what would be deleted without making changes
poweradmin zones delete --name example.com --dry-run
poweradmin zones delete --interactive --dry-run

# Export a zone as BIND zone file
poweradmin zones export --name example.com
poweradmin zones export --name example.com > example.com.zone

# Import a zone from a BIND zone file
poweradmin zones import --file example.com.zone
poweradmin zones import --file example.com.zone --zone-name example.com
poweradmin zones import --file example.com.zone --create-zone
poweradmin zones import --file example.com.zone --dry-run

# List all metadata entries for a zone
poweradmin zones metadata --name example.com
poweradmin zones metadata --name example.com -o json

# Get metadata for a specific kind
poweradmin zones metadata-get --name example.com --kind ALLOW-AXFR-FROM

# Set (replace) metadata values for a kind
poweradmin zones metadata-set --name example.com \
  --kind ALLOW-AXFR-FROM \
  --values 192.0.2.10,AUTO-NS

# Delete metadata for a kind
poweradmin zones metadata-delete --name example.com --kind ALLOW-AXFR-FROM --yes
```

### Records

```bash
# List all records in a zone
poweradmin records list --zone-name example.com
poweradmin records list --zone-name example.com --type A
poweradmin records list --zone-name example.com --sort type
poweradmin records list --zone-name example.com -o json

# Get a record by ID
poweradmin records get --zone-name example.com --id <record-id>

# Create a record
poweradmin records create \
  --zone-name example.com \
  --name www.example.com \
  --type A \
  --content 1.2.3.4 \
  --ttl 3600

# Create a record interactively
poweradmin records create --interactive

# Create and capture the ID
RECORD_ID=$(poweradmin records create \
  --zone-name example.com \
  --name www.example.com \
  --type A \
  --content 1.2.3.4 -q)

# Update a record
poweradmin records update \
  --zone-name example.com \
  --id <record-id> \
  --content 5.6.7.8

# Delete a record
poweradmin records delete \
  --zone-name example.com \
  --id <record-id> --yes

# Delete records interactively (pick a zone, then multi-select records)
poweradmin records delete --interactive

# Preview what would be deleted without making changes
poweradmin records delete --zone-name example.com --id <record-id> --dry-run
```

### Users

```bash
# List all users
poweradmin users list
poweradmin users list --active
poweradmin users list --active=false
poweradmin users list --sort username

# Get a user by name or ID
poweradmin users get --name patrick
poweradmin users get --id 1

# Create a user (password will be prompted)
poweradmin users create \
  --username patrick \
  --email patrick@example.com \
  --fullname "Patrick Omland"

# Create a user interactively
poweradmin users create --interactive

# Create a user with password via flag (not recommended)
poweradmin users create \
  --username patrick \
  --password secret123 \
  --email patrick@example.com

# Update a user
poweradmin users update --name patrick --email new@example.com
poweradmin users update --id 1 --active=false

# Delete a user
poweradmin users delete --name patrick --yes

# Delete a user who still owns zones, handing them over to another user
poweradmin users delete --name patrick --transfer-to anna --yes

# Delete users interactively (multi-select)
poweradmin users delete --interactive

# Preview what would be deleted without making changes
poweradmin users delete --name patrick --dry-run

# Assign a permission template
poweradmin users set-permission-template --name patrick --template-id 2
```

### Groups

```bash
# List all groups
poweradmin groups list
poweradmin groups list -o json

# Get a group by name or ID
poweradmin groups get --name Administrators
poweradmin groups get --id 1

# Create a group
poweradmin groups create --name "Zone Editors" --description "Can edit zone records"
GROUP_ID=$(poweradmin groups create --name "Zone Editors" -q)

# Create a group interactively
poweradmin groups create --interactive

# Update a group
poweradmin groups update --name "Zone Editors" --new-name "DNS Editors"

# Delete a group
poweradmin groups delete --name "DNS Editors" --yes

# Delete groups interactively (multi-select)
poweradmin groups delete --interactive

# Preview what would be deleted without making changes
poweradmin groups delete --name "Zone Editors" --dry-run

# Manage members
poweradmin groups members --name Administrators
poweradmin groups member-add --group-id 1 --user-id 2
poweradmin groups member-remove --group-id 1 --user-id 2

# Manage zones
poweradmin groups zones --name Administrators
poweradmin groups zone-add --group-id 1 --zone-id 78
poweradmin groups zone-remove --group-id 1 --zone-id 78
```

### Permission Templates

```bash
# List all permission templates
poweradmin permission-templates list
poweradmin permission-templates list -o json

# Get a template by name or ID (shows full permission list)
poweradmin permission-templates get --name Administrator
poweradmin permission-templates get --id 1 -o json

# Create a template
poweradmin permission-templates create --name "Zone Editors" --description "Can edit zone records"
poweradmin permission-templates create --name "Zone Editors" --type group --permissions 1,2,3

# Create a template interactively
poweradmin permission-templates create --interactive

# Create and capture the ID
PT_ID=$(poweradmin permission-templates create --name "Zone Editors" -q)

# Update a template
poweradmin permission-templates update --name "Zone Editors" --new-name "DNS Editors"
poweradmin permission-templates update --name "Zone Editors" --permissions 1,2,3,4

# Delete a template
poweradmin permission-templates delete --name "Zone Editors" --yes

# Delete templates interactively (multi-select)
poweradmin permission-templates delete --interactive

# Preview what would be deleted without making changes
poweradmin permission-templates delete --name "Zone Editors" --dry-run
```

### Version

```bash
poweradmin version
```

## Output Formats

| Flag | Description |
|------|-------------|
| `-o table` | Human-readable aligned table, long content truncated (default) |
| `-o full` | Human-readable aligned table, full content |
| `-o json` | JSON output, suitable for scripting and piping into `jq` |
| `-o yaml` | YAML output, suitable for piping into tools like `yq` |

Color output is automatically enabled when stdout is a terminal and suppressed
when piping or redirecting output.

## Scripting

```bash
# Get all zone IDs
poweradmin zones list -o json | jq '.zones[].id'

# Create a zone and immediately add a record
ZONE_ID=$(poweradmin zones create example.com -q)
poweradmin records create \
  --zone-id $ZONE_ID \
  --name www.example.com \
  --type A \
  --content 1.2.3.4

# Delete all records of a specific type
poweradmin records list --zone-name example.com --type TXT -o json \
  | jq -r '.records[].id' \
  | xargs -I{} poweradmin records delete --zone-name example.com --id {} --yes

# Backup all zones as BIND zone files
poweradmin zones list -o json | jq -r '.zones[].name' | \
  xargs -I{} sh -c 'poweradmin zones export --name {} > {}.zone'

# Migrate a zone to a new Poweradmin instance
poweradmin zones export --name example.com > example.com.zone
POWERADMIN_URL=https://new-dns.example.com \
  poweradmin zones import --file example.com.zone --create-zone
```

## Security

Docker images are signed with [cosign](https://github.com/sigstore/cosign) using keyless signing via GitHub OIDC. Verify an image with:

```bash
cosign verify \
  --certificate-identity-regexp="https://github.com/contentways/poweradmin-cli" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com" \
  contentwaysorg/poweradmin-cli:latest
```

Each release also includes a Software Bill of Materials (SBOM) generated by [Syft](https://github.com/anchore/syft).

## Documentation

Full reference documentation is available in [docs/reference](docs/reference/).

## License

MIT — Copyright (c) 2026 Contentways
