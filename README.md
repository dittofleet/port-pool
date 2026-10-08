<img src="assets/icon.svg" width="80" alt="port-pool icon">

# port-pool

Simple port allocation tool for testing multiple dev environments in parallel. Useful for:

- Git worktrees that each need unique ports
- Running multiple copies of the same repository
- Running multiple projects simultaneously without port conflicts

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/dittofleet/.github/main/install.sh | sh -s port-pool
```

Installs the latest release to `~/.local/bin/port-pool`, with a starter config listing ports to stay off. Override the install location with `PORT_POOL_INSTALL_DIR`:

```sh
curl -fsSL https://raw.githubusercontent.com/dittofleet/.github/main/install.sh | PORT_POOL_INSTALL_DIR=/usr/local/bin sh -s port-pool
```

Supported platforms: macOS (arm64, x64), Linux (arm64, x64).

## Updating

`port-pool` checks once per day for new releases and prints a hint to stderr when an update is available. Run `port-pool update` to upgrade in place.

The check is automatically skipped when:

- `CI` is set
- `PORT_POOL_NO_UPDATE_CHECK` is set
- stderr is not a TTY (piped or redirected)
- the binary was built locally (version is `dev`)

## Uninstalling

```bash
port-pool uninstall          # prompts for confirmation
port-pool uninstall --yes    # skip prompt
```

Removes the binary, `~/.config/port-pool/` (global config and any per-pool configs), and `~/.local/share/port-pool/` (state and update cache). `.env` files in project directories are not touched.

## Setup

The install script creates a starter config at `~/.config/port-pool/config.json` that uses the default range of 3000 to 9999 and has common dev-tool ports (databases, Vite, Django, etc.) excluded. Edit it to suit your machine (see [Configuration](#configuration)).

## Usage

```bash
# Allocate ports for a project (updates .env)
port-pool provision /path/to/project

# Release ports when done
port-pool release /path/to/project

# Idempotent: provision if needed, repair .env drift if not. Designed to be
# chained inline before a project's dev command.
port-pool ensure /path/to/project

# See all current allocations
port-pool list

# Remove allocations whose project directory or config is gone
port-pool prune              # actually removes
port-pool prune --dry-run    # just lists what would be removed
```

### Recommended: chain `ensure` into your dev command

Wire `port-pool ensure .` directly into the project's dev script so a fresh clone "just works":

```json
{
  "scripts": {
    "dev": "port-pool ensure . && vite"
  }
}
```

On the first `npm run dev` after cloning, ports are auto-provisioned and the configured env files are created. Subsequent runs are silent fast no-ops. If an env file has been hand-edited away from the templated values, `ensure` repairs it before the dev server starts. If the project holds a port below 1024, above 65535, or one that browsers block (see [Configuration](#configuration)), `ensure` replaces its ports with a new block.

For CI or pre-commit checks, use `--check` to fail without modifying anything:

```bash
port-pool ensure --check .
```

`--check` exits with status 1 if the project is not provisioned, holds a port below 1024, above 65535, or one that browsers block, or has an env file that drifted. Otherwise it's silent.

Multiple parallel `ensure` calls (e.g. across worktrees starting their dev servers simultaneously) are safe; allocations serialize via a sidecar lock at `~/.local/share/port-pool/state.json.lock`.

## Configuration

`port-pool` uses two configs:

### Global pool config (`~/.config/port-pool/config.json`)

Defines the machine-level port pool:

```json
{
  "schemaVersion": 1,
  "portRangeStart": 3000,
  "portRangeEnd": 9999,
  "excludedPorts": [
    3000, 3001, 3306,
    4000, 4200,
    5000, 5173, 5432, 5500,
    6379,
    8000, 8080, 8081, 8443, 8888,
    9000, 9090, 9200
  ]
}
```

- `schemaVersion`: required integer. Identifies the file's schema. The current schema is `1`. If a future port-pool release changes the schema, this number bumps and the binary errors clearly when it sees an old (or unknown) version.
- `portRangeStart`/`portRangeEnd`: inclusive bounds of the pool. Optional. A missing bound defaults to 3000 (start) or 9999 (end). A range with no ports between 1024 and 65535 is rejected.
- `excludedPorts`: ports inside the range that should never be allocated (typically defaults of common dev tools — Postgres, Redis, Vite, Django, Jupyter, etc.).

Whatever the range, port-pool never hands out ports below 1024 or above 65535, or ports that browsers refuse to load (the Fetch standard's [bad ports](https://fetch.spec.whatwg.org/#port-blocking), such as 6000 and 6566). None of these need to be listed. Existing allocations that hold one are replaced by `ensure`.

Respects `$XDG_CONFIG_HOME` if set.

### Per-project config (`port-pool.config.json`)

Lives in each project's root, alongside `package.json` / `tsconfig.json`. Typically committed to the project repo (unlike resolved env files, which contain the port numbers and are usually git-ignored).

```json
{
  "schemaVersion": 1,
  "portNames": ["web", "api", "db"],
  "envFiles": {
    ".env": {
      "DB_PORT": "${db}"
    },
    "apps/web/.env.local": {
      "PORT":    "${web}",
      "API_URL": "http://localhost:${api}",
      "WS_URL":  "ws://localhost:${api}"
    },
    "apps/api/.env": {
      "PORT":         "${api}",
      "DATABASE_URL": "postgres://localhost:${db}/app"
    }
  }
}
```

- `schemaVersion`: required integer. Identifies the file's schema. The current schema is `1`. Future port-pool releases that change the schema bump this number; the binary errors clearly when it sees an old (or unknown) version.
- `portNames`: project-global list of port names. The length determines how many ports are allocated for the project. Names must match `[a-zA-Z_][a-zA-Z0-9_]*` and be unique. Templates reference ports as `${NAME}`.
- `envFiles`: map of relative env file path to env-var templates. Each value is a record of `ENV_VAR_NAME` → template string. Same `${NAME}` can appear in multiple files and resolves to the same port (that's the point).

Path rules: env file paths are relative to the project root, may be in subdirectories (e.g. `apps/web/.env`), and reject `..` segments and absolute paths. JSON keys enforce uniqueness.

When provisioning:

- **Keys not mentioned in `port-pool.config.json` are left untouched** — port-pool only manages the keys you list, so existing secrets and unrelated settings in the same file are safe.
- Existing managed keys are updated in place; missing managed keys are appended.
- Files are created if they don't exist.

## State

Allocations are stored at `~/.local/share/port-pool/state.json` (respects `$XDG_DATA_HOME`).

## Development

Running from source still works:

```bash
go run . ensure /path/to/project
```

To isolate dev runs from your real allocations, override the XDG variables:

```bash
XDG_CONFIG_HOME=/tmp/dev/config XDG_DATA_HOME=/tmp/dev/data go run . ensure /tmp/test-project
```

### Building from source

```sh
go build -o dist/port-pool .
cp dist/port-pool ~/.local/bin/
```

Source builds report version `dev`, which disables the auto-update check.
