# Changelog

## [v2.0.1](https://github.com/contentways/poweradmin-cli/releases/tag/v2.0.1)

### Bug Fixes

- derive version/commit from build info when unset via ldflags

## [v2.0.0](https://github.com/contentways/poweradmin-cli/releases/tag/v2.0.0)

### Features

- add --ttl flag for NS records in zones create
- **BREAKING**: migrate module path to github.com/contentways/poweradmin-go/v3

### Bug Fixes

- import sorting
- correct repository case in releaser-pleaser workflow

## [v1.3.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.3.0)

### Features

- add --ttl flag for NS records in zones create

## [v1.2.1](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.2.1)

### Bug Fixes

- use GroupNameCompletion for groups commands, add PermissionTemplateNameCompletion

## [v1.2.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.2.0)

### Features

- add zones export command for BIND zone file format
- add zones import command with BIND zone file parser
- update README.md
- add permission-templates commands (list, get, create, update, delete)

## [v1.3.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.3.0)

### Features

- switch to keyless cosign signing via GitHub OIDC

## [v1.2.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.2.0)

### Features

- add cosign image signing to release workflow

## [v1.1.3](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.1.3)

### Bug Fixes

- correct Docker Hub org name to contentwaysorg

## [v1.1.2](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.1.2)

### Bug Fixes

- simplify Dockerfile - use pre-built binary from GoReleaser context

## [v1.2.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.2.0)

### Features

- migrate to dockers_v2, update Dockerfile for multi-arch buildx

### Bug Fixes

- use StringSlice for --nameserver flag to support comma-separated values

## [v1.1.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.1.0)

### Features

- add delete confirmation prompt with --yes/-y flag
- add --no-header flag to list commands, --quiet/-q to create/delete
- add color output for TTY terminals
- add records get command with client-side ID lookup
- add sort flags, records get, records/users filter, stderr error output, base package refactor
- add shell completion via API for zone, user and group names

## [v0.4.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v0.4.0)

### Features

- add version command with build-time version injection
- add records update command with tests
- add --type and --name-filter flags to zones list
- add record update, version command, zone filters, fix tests

## [v0.3.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v0.3.0)

### Features

- add groups commands (list, get, create, update, delete, members, zones)

### Bug Fixes

- update tests for new schema package JSON output

## [v0.2.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v0.2.0)

### Features

- add short flags -u (--url), -k (--api-key), -o (--output)
- add schema package for clean JSON output, add root objects and omit empty fields
