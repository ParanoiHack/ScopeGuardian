# Scope-Guardian

Scope-Guardian is a CLI tool that runs security scanners on your codebase and synchronises the results with [DefectDojo](https://github.com/DefectDojo/django-DefectDojo). It can optionally enforce a security gate that blocks a CI/CD pipeline when finding counts exceed configurable thresholds.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [CLI Usage](#cli-usage)
4. [Configuration File (`config.toml`)](#configuration-file-configtoml)
5. [Environment Variables](#environment-variables)
6. [How Engagements Are Handled](#how-engagements-are-handled)
7. [How the Sync Feature Works](#how-the-sync-feature-works)
8. [How the Security Gate Works](#how-the-security-gate-works)
9. [Running with Docker](#running-with-docker)
10. [Local DefectDojo Setup](#local-defectdojo-setup)

---

## Prerequisites

- Go 1.25+ (only needed when building from source)
- [KICS](https://github.com/Checkmarx/kics) binary available at `/opt/kics/bin/kics` (pre-installed in the Docker image)
- [Syft](https://github.com/anchore/syft) binary available at `/opt/syft/bin/syft` (pre-installed in the Docker image; required when `[grype]` is configured)
- [Grype](https://github.com/anchore/grype) binary available at `/opt/grype/bin/grype` (pre-installed in the Docker image; required when `[grype]` is configured)
- A running DefectDojo instance and an API access token (required only when `--sync` is used)

---

## Quick Start

```bash
# Build the binary
go build -o scope-guardian .

# Run a basic scan (no sync, no gate)
SCAN_DIR=/path/to/repos ./scope-guardian \
  --projectName my-service \
  --branch main \
  ./config.toml

# Run a scan, sync results to DefectDojo, and enforce a security gate
SCAN_DIR=/path/to/repos \
DD_URL=http://localhost:8080 \
DD_ACCESS_TOKEN=<your-token> \
./scope-guardian \
  --projectName my-service \
  --branch main \
  --sync \
  --threshold critical=1,high=5 \
  ./config.toml
```

---

## CLI Usage

```
scope-guardian [flags] <config-file>
```

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--projectName` | string | yes | Name of the project being scanned. Must match the product name in DefectDojo when `--sync` is used. |
| `--branch` | string | yes | Branch being scanned (e.g. `main`, `feature/my-branch`). |
| `--sync` | bool | no | Upload scan results to DefectDojo. Requires `DD_URL` and `DD_ACCESS_TOKEN`. Default: `false`. |
| `--threshold` | string | no | Comma-separated severity thresholds that define the security gate (see [Security Gate](#how-the-security-gate-works)). |
| `<config-file>` | path | yes | Path to the TOML configuration file. |

### Execution Order

```
Parse flags → Load config.toml → Initialize scanners
  → Phase 1: Run prerequisite scanners concurrently (Syft SBOM generation)
  → Phase 2: Run dependent/independent scanners concurrently (Grype, KICS)
           Any scanner whose prerequisite failed is skipped automatically
  → Load findings → [Sync to DefectDojo] → Display findings
  → [Evaluate security gate → exit(-1) on failure]
```

When both `--sync` and `--threshold` are provided the gate is evaluated against the findings already stored in DefectDojo (post-deduplication) rather than the raw local scan output.

---

## Configuration File (`config.toml`)

The configuration file is a [TOML](https://toml.io) document that controls which scanners run and how engagements are managed.

```toml
title = "Scope-guardian configuration file"   # Optional human-readable label

# Branches whose DefectDojo engagements are given a one-year end date.
# All other branches receive a one-week end date.
protected_branches = ["main", "master"]

# KICS – infrastructure-as-code scanner
[kics]
# Directory to scan, relative to the SCAN_DIR environment variable.
path = "./my-service"
# Infrastructure platform type. Passed as --type to KICS.
# Examples: "Dockerfile", "Terraform", "CloudFormation", "Kubernetes", "Ansible"
platform = "Dockerfile"

# Grype – software-composition analysis (SCA) vulnerability scanner.
# Enabling this section also enables Syft SBOM generation as a prerequisite.
[grype]
# Comma-separated vulnerability states to ignore.
# Common values: "not-fixed", "unknown", "wont-fix"
ignore_states = "not-fixed,unknown,wont-fix"
# When true, Syft resolves transitive Java dependencies from Maven Central.
# This increases scan accuracy for Java projects but significantly increases scan time.
transitive_libraries = false
# Optional list of path patterns to exclude from Grype scanning.
# exclude = ["**/vendor/**", "**/testdata/**"]
```

### Fields Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | no | Human-readable label; not used programmatically. |
| `protected_branches` | string array | no | Branches whose engagements get a 1-year end date. Defaults to empty (all branches get 1 week). |
| `[kics].path` | string | yes* | Path to the directory to scan. Resolved as `$SCAN_DIR/<path>`. |
| `[kics].platform` | string | no | KICS platform filter (e.g. `Dockerfile`). When omitted KICS scans all supported types. |
| `[grype].ignore_states` | string | no | Comma-separated Grype vulnerability states to suppress (e.g. `not-fixed,unknown,wont-fix`). |
| `[grype].transitive_libraries` | bool | no | When `true`, Syft resolves transitive Java dependencies via Maven Central. Default: `false`. |
| `[grype].exclude` | string array | no | Path glob patterns to exclude from Grype scanning (e.g. `["**/vendor/**"]`). |

\* Required only if you want KICS scanning to run. Omitting the entire `[kics]` section disables the scanner.

Omitting the entire `[grype]` section disables both Grype and the Syft SBOM generation step.

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SCAN_DIR` | yes | Base directory for scan operations. Scan paths and result files are resolved relative to this value. |
| `DD_URL` | when `--sync` | Base URL of the DefectDojo instance (e.g. `http://localhost:8080`). |
| `DD_ACCESS_TOKEN` | when `--sync` | DefectDojo API token. Generate one in DefectDojo under **User → API v2 Key**. |

Copy `.env.example` to `.env` and fill in the values for local development:

```bash
cp .env.example .env
```

---

## How Engagements Are Handled

Scope-Guardian uses a single DefectDojo **engagement** per project/branch combination to store all findings for that branch. Engagements are managed automatically — you never have to create or update them manually.

### Engagement Naming

Every engagement is named `<projectName>-<branch>`, for example:

- `my-service-main`
- `my-service-feature-my-branch`

The project name must correspond to an existing **Product** in DefectDojo with an exact name match.

### Engagement Lifecycle

When `--sync` is used the following logic runs on every invocation:

1. **Look up the product** in DefectDojo by exact name (`projectName`).
2. **List all engagements** for that product (all pages).
3. **Search for an engagement** whose name matches `<projectName>-<branch>`.
   - **Found, end date still valid** → reuse it as-is.
   - **Found, end date in the past** → automatically extend the end date and reuse the engagement.
   - **Not found** → create a new engagement.
4. **Upload scan results** into the engagement.

### Engagement Duration

The end date is determined by whether the branch appears in `protected_branches`:

| Branch type | End date |
|------------|----------|
| Protected (e.g. `main`, `master`) | 1 year from today |
| Unprotected (feature branches, etc.) | 1 week from today |

### Engagement Details

New engagements are created with the following attributes:

| Attribute | Value |
|-----------|-------|
| Type | CI/CD |
| Status | In Progress |
| Tags | `SCOPE-GUARDIAN`, `<branch>` |
| Deduplication on engagement | disabled |

---

## How the Sync Feature Works

Passing `--sync` on the command line uploads the scan results from every registered scanner to DefectDojo after the scan completes.

### Step-by-step

1. A DefectDojo service client is created using `DD_URL` and `DD_ACCESS_TOKEN`.
2. The engagement ID is resolved (see [Engagements](#how-engagements-are-handled)).
3. For each registered scanner, the scanner's `Sync` method is called with the engagement ID and branch.

### KICS Sync Behaviour

The KICS scanner uploads its JSON output file to DefectDojo via the `/api/v2/import-scan/` endpoint as a `multipart/form-data` request. The following options are set on every import:

| Option | Value | Effect |
|--------|-------|--------|
| Scan type | `KICS Scan` | Tells DefectDojo which parser to use |
| Severity threshold | `Info` | Import findings of all severities |
| Group by | `finding_title` | Merge findings with the same title |
| Create finding groups | `true` | Group related findings together |
| Apply tags to findings | `true` | Tag each finding with `IACST` |
| Close old findings | `true` | Findings absent from the new scan are closed automatically |
| Branch tag | `<branch>` | Associates the results with the scanned branch |

### Grype Sync Behaviour

The Grype scanner uploads its JSON output file to DefectDojo via the `/api/v2/import-scan/` endpoint as a `multipart/form-data` request. The following options are set on every import:

| Option | Value | Effect |
|--------|-------|--------|
| Scan type | `Anchore Grype` | Tells DefectDojo which parser to use |
| Severity threshold | `Info` | Import findings of all severities |
| Group by | `finding_title` | Merge findings with the same title |
| Create finding groups | `true` | Group related findings together |
| Apply tags to findings | `true` | Tag each finding with `SCA` |
| Close old findings | `true` | Findings absent from the new scan are closed automatically |
| Branch tag | `<branch>` | Associates the results with the scanned branch |

### Security Gate with Sync

When both `--sync` and `--threshold` are set, the gate is evaluated against **DefectDojo's deduplicated findings** rather than the raw local scan output. This means:

- Duplicate or previously-closed findings do not inflate the count.
- Only active findings surviving DefectDojo's deduplication logic are counted.

---

## How the Security Gate Works

The security gate fails the pipeline (exit code `-1`) when the number of findings at or above a configured severity level meets or exceeds the configured limit.

### Threshold Syntax

```
--threshold <severity>=<count>[,<severity>=<count>...]
```

Supported severity values (case-insensitive): `critical`, `high`, `medium`, `low`, `info`.

```bash
# Fail on any critical finding
--threshold critical=1

# Fail on 1+ critical OR 5+ high findings
--threshold critical=1,high=5

# Fail on 10+ medium-or-above findings
--threshold medium=10
```

### Evaluation Logic

For each threshold rule:

1. Count findings whose severity is **equal to or higher than** the threshold severity.

   Severity ranking (highest to lowest): `CRITICAL` > `HIGH` > `MEDIUM` > `LOW` > `INFO`

2. If the count is **≥ the configured value**, the gate **fails** and the process exits with code `-1`.
3. All threshold rules must pass for the gate to pass.

### Finding Source

| Flags used | Findings evaluated |
|-----------|-------------------|
| `--threshold` only | Local scan output from all scanners |
| `--threshold` + `--sync` | Active findings fetched from DefectDojo |

---

## Running with Docker

The provided `Dockerfile` builds a multi-stage image that bundles Scope-Guardian together with KICS, OpenGrep, Grype, and Syft.

```bash
# Build the image
docker build -t scope-guardian .

# Run a scan
docker run --rm \
  -v /path/to/your/project:/tmp/data/project \
  -v /path/to/config.toml:/config.toml \
  -e SCAN_DIR=/tmp/data \
  -e DD_URL=http://host.docker.internal:8080 \
  -e DD_ACCESS_TOKEN=<your-token> \
  scope-guardian \
  --projectName my-service \
  --branch main \
  --sync \
  /config.toml
```

Inside the container `SCAN_DIR` defaults to `/tmp/data`.

---

## Local DefectDojo Setup

A `docker-compose.yml` is provided to spin up a local DefectDojo instance backed by PostgreSQL and Redis.

```bash
# 1. Configure credentials
cp .env.example .env
# Edit .env – change passwords and set a strong DD_SECRET_KEY

# 2. Start DefectDojo (this takes a minute on first run)
docker compose up -d

# 3. Open DefectDojo in your browser
open http://localhost:8080

# 4. Log in with the admin credentials you set in .env
#    Default: admin / changeme

# 5. Generate an API token
#    Profile → API v2 Key → copy the token

# 6. Set the token in your environment
export DD_URL=http://localhost:8080
export DD_ACCESS_TOKEN=<your-token>
```

### Creating a DefectDojo Product

Before running Scope-Guardian with `--sync`, create a **Product** in DefectDojo whose name matches the `--projectName` value you will pass on the command line.

1. In DefectDojo go to **Products → Add Product**.
2. Set **Name** to exactly the value you will pass as `--projectName` (e.g. `my-service`).
3. Save. Scope-Guardian will manage engagements under this product automatically.
