# poweradmin-cli

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
go install github.com/contentways/poweradmin-cli/v2@latest
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
| `--output` | `-o` | Output format: `table`, `full`, `json` |
| `--no-header` | | Suppress table header row |
| `--quiet` | `-q` | Only print ID (create) or suppress output (delete) |
| `--yes` | `-y` | Skip delete confirmation prompt |

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

# Create and capture the ID
ZONE_ID=$(poweradmin zones create example.com -q)

# Delete a zone
poweradmin zones delete --name example.com
poweradmin zones delete --name example.com --yes

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

# Update a group
poweradmin groups update --name "Zone Editors" --new-name "DNS Editors"

# Delete a group
poweradmin groups delete --name "DNS Editors" --yes

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

# Create and capture the ID
PT_ID=$(poweradmin permission-templates create --name "Zone Editors" -q)

# Update a template
poweradmin permission-templates update --name "Zone Editors" --new-name "DNS Editors"
poweradmin permission-templates update --name "Zone Editors" --permissions 1,2,3,4

# Delete a template
poweradmin permission-templates delete --name "Zone Editors" --yes
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
