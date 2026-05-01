# ScopeGuardian – Context for AI Assistants

## Project Overview

**ScopeGuardian** is a Go CLI tool that orchestrates security scanners (KICS, Grype/Syft, OpenGrep) against a codebase and optionally synchronises findings with [DefectDojo](https://github.com/DefectDojo/django-DefectDojo). It can enforce a security gate that exits with `-1` to block CI/CD pipelines when configured severity thresholds are exceeded.

## Repository Layout

```
ScopeGuardian/
├── main.go                        # Entry point: parse → load → scan → display → gate
├── config.toml                    # Sample configuration file
├── go.mod / go.sum                # Go module (module name: ScopeGuardian, go 1.24.4)
├── Dockerfile                     # Multi-stage image (bundles KICS, Grype, Syft, OpenGrep)
├── docker-compose.yml             # Local DefectDojo stack (PostgreSQL + Redis + Nginx)
├── .env.example                   # Template for environment variables
├── connectors/
│   └── defectdojo/                # DefectDojo API client (import/reimport, findings, engagements)
│       └── client/                # HTTP client wrapper
├── display/                       # Banner, findings table, JSON/CSV/raw dump
├── domains/
│   ├── interfaces/                # ScanServiceImpl interface
│   └── models/                    # Finding model, FindingStatus, FilterFindingsByStatus
├── engine/
│   ├── engine.go                  # Two-phase parallel runner (prerequisites → dependents)
│   ├── engine_test.go
│   └── const.go                   # Log message constants
├── environnement_variable/        # Loads SCAN_DIR, DD_URL, DD_ACCESS_TOKEN from env
├── exec/                          # Shell command execution helpers
├── features/
│   ├── scans/
│   │   ├── kics/                  # KICS IaC scanner integration
│   │   ├── grype/                 # Grype SCA scanner integration
│   │   ├── opengrep/              # OpenGrep SAST scanner integration
│   │   └── syft/                  # Syft SBOM generator (prerequisite for Grype)
│   ├── security-gate/             # Threshold evaluation logic
│   └── sync/                      # DefectDojo engagement resolution and finding marking
├── loader/                        # TOML config loader
├── logger/                        # Structured logging (slog wrapper + null logger)
└── parser/                        # CLI flag parser (--projectName, --branch, --sync, etc.)
```

## Common Commands

```bash
# Build
go build -o ScopeGuardian .

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a basic scan (no sync, no gate)
SCAN_DIR=/path/to/repos ./ScopeGuardian --projectName my-service --branch main ./config.toml

# Run a scan and write findings to JSON
SCAN_DIR=/path/to/repos ./ScopeGuardian --projectName my-service --branch main -o /tmp/findings.json ./config.toml

# Run a scan, sync to DefectDojo, enforce gate
SCAN_DIR=/path/to/repos DD_URL=http://localhost:8080 DD_ACCESS_TOKEN=<token> \
  ./ScopeGuardian --projectName my-service --branch main --sync --threshold critical=1,high=5 ./config.toml

# Build Docker image
docker build -t ScopeGuardian .
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SCAN_DIR` | always | Base directory; scan paths and result files are resolved relative to this |
| `DD_URL` | with `--sync` | Base URL of the DefectDojo instance (e.g. `http://localhost:8080`) |
| `DD_ACCESS_TOKEN` | with `--sync` | DefectDojo API v2 token |

## Key CLI Flags

| Flag | Description |
|------|-------------|
| `--projectName` | Project/product name; must match a DefectDojo Product name when using `--sync` |
| `--branch` | Branch being scanned |
| `--sync` | Upload results to DefectDojo and fetch back statuses |
| `--threshold` | `severity=count[,...]` — security gate thresholds |
| `--filter` | Comma-separated statuses to display/write (`ACTIVE`, `INACTIVE`, `DUPLICATE`). Default: `ACTIVE` |
| `-q` | Quiet mode (suppress logs) |
| `-o` | Output file path for findings |
| `--format` | Output format: `json` (default), `csv`, `raw` |

## Architecture Notes

### Engine (two-phase parallel scan)
1. **Prerequisites** — Syft SBOM generation runs first; failure is recorded.
2. **Dependent/independent scanners** — Grype, KICS, OpenGrep run concurrently. Grype is skipped if Syft failed.

### Scanner interface (`domains/interfaces/ScanServiceImpl`)
Every scanner implements:
- `Start() (bool, error)` — runs the scanner binary, writes output file to `$SCAN_DIR/results/`
- `LoadFindings() ([]models.Finding, error)` — parses the output file into `Finding` structs
- `Sync(engagementId int, branch string, ddService) error` — uploads results to DefectDojo

### Finding model (`domains/models/Finding`)
Fields include: `Title`, `Severity`, `Description`, `FilePath`, `Status` (`ACTIVE`/`INACTIVE`/`DUPLICATE`).

### DefectDojo sync behaviour
- **First run** → `POST /api/v2/import-scan/`
- **Subsequent runs** → `POST /api/v2/reimport-scan/` (closes old findings, preserves triage decisions)
- Engagement name pattern: `<projectName>-<branch>`
- Protected branches → 1-year engagement end date; others → 1-week end date

### Security gate
- Counts findings **≥ configured severity** (CRITICAL > HIGH > MEDIUM > LOW > INFO).
- When `--sync` is active the gate uses DefectDojo's deduplicated ACTIVE findings, not raw local output.
- Exit code `-1` on failure.

## Configuration File (`config.toml`)

```toml
title = "Scope-guardian configuration file"
protected_branches = ["main", "master"]
path = "./my-service"          # Relative to SCAN_DIR

[kics]
platform = "Dockerfile"        # KICS --type filter

[grype]
ignore_states = "not-fixed,unknown,wont-fix"
transitive_libraries = false   # true = resolve transitive Java deps (slow)

[opengrep]
exclude = ["tests/**"]
exclude_rule = []

# [proxy]
# http_proxy  = "http://proxy.company.com:3128"
# https_proxy = "http://proxy.company.com:3128"
# no_proxy    = "localhost,127.0.0.1"
# ssl_cert_file = "/path/to/ca.pem"
```

## Code Conventions

- **Module path**: `ScopeGuardian` (no domain prefix)
- **Logging**: use `logger.Info`, `logger.Error`, `logger.Err(err)` from the `logger` package — never `fmt.Println` or `log.*`
- **Log constants**: define log message strings as `const` in the package (see `engine/const.go`)
- **Tests**: use `github.com/stretchr/testify` and `github.com/golang/mock`
- **Error handling**: log and continue (scanner errors do not abort the whole run)
- **No global state**: scanners are instantiated via `Get<Scanner>Service(config)` factory functions
- **Proxy forwarding**: pass proxy env vars to all scanner sub-processes; set both uppercase and lowercase variants

## External Tool Binary Paths (inside Docker / expected locations)

| Tool | Path |
|------|------|
| KICS | `/opt/kics/bin/kics` |
| OpenGrep | `/opt/opengrep/bin/opengrep` |
| Syft | `/opt/syft/bin/syft` |
| Grype | `/opt/grype/bin/grype` |

## Testing

- Unit tests live alongside source files (`*_test.go`).
- Run `go test ./...` from the repo root.
- Mocks are generated with `github.com/golang/mock`.
- No integration test suite is bundled; integration testing requires a live DefectDojo instance.
