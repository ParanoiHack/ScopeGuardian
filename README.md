# ScopeGuardian

ScopeGuardian is a CLI tool that runs security scanners on your codebase and synchronises the results with [DefectDojo](https://github.com/DefectDojo/django-DefectDojo). It can optionally enforce a security gate that blocks a CI/CD pipeline when finding counts exceed configurable thresholds.

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
- [OpenGrep](https://github.com/opengrep/opengrep) binary available at `/opt/opengrep/bin/opengrep` (pre-installed in the Docker image; required when `[opengrep]` is configured)
- [Syft](https://github.com/anchore/syft) binary available at `/opt/syft/bin/syft` (pre-installed in the Docker image; required when `[grype]` is configured)
- [Grype](https://github.com/anchore/grype) binary available at `/opt/grype/bin/grype` (pre-installed in the Docker image; required when `[grype]` is configured)
- A running DefectDojo instance and an API access token (required only when `--sync` is used)

---

## Quick Start

```bash
# Build the binary
go build -o ScopeGuardian .

# Run a basic scan (no sync, no gate)
SCAN_DIR=/path/to/repos ./ScopeGuardian \
  --projectName my-service \
  --branch main \
  ./config.toml

# Run a scan with quiet mode (no logs)
SCAN_DIR=/path/to/repos ./ScopeGuardian \
  --projectName my-service \
  --branch main \
  -q \
  ./config.toml

# Run a scan and write findings to a file (JSON by default)
SCAN_DIR=/path/to/repos ./ScopeGuardian \
  --projectName my-service \
  --branch main \
  -o /tmp/findings.json \
  ./config.toml

# Run a scan and write findings as CSV
SCAN_DIR=/path/to/repos ./ScopeGuardian \
  --projectName my-service \
  --branch main \
  -o /tmp/findings.csv \
  --format csv \
  ./config.toml

# Run a scan showing active and duplicate findings
SCAN_DIR=/path/to/repos ./ScopeGuardian \
  --projectName my-service \
  --branch main \
  --filter ACTIVE,DUPLICATE \
  ./config.toml

# Run a scan, sync results to DefectDojo, and enforce a security gate
SCAN_DIR=/path/to/repos \
DD_URL=http://localhost:8080 \
DD_ACCESS_TOKEN=<your-token> \
./ScopeGuardian \
  --projectName my-service \
  --branch main \
  --sync \
  --threshold critical=1,high=5 \
  ./config.toml
```

---

## CLI Usage

```
ScopeGuardian [flags] <config-file>
```

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--projectName` | string | yes | Name of the project being scanned. Must match the product name in DefectDojo when `--sync` is used. |
| `--branch` | string | yes | Branch being scanned (e.g. `main`, `feature/my-branch`). |
| `--sync` | bool | no | Upload scan results to DefectDojo. Requires `DD_URL` and `DD_ACCESS_TOKEN`. Default: `false`. |
| `--threshold` | string | no | Comma-separated severity thresholds that define the security gate (see [Security Gate](#how-the-security-gate-works)). |
| `--filter` | string | no | Comma-separated finding statuses to include in display and file output. Accepted values: `ACTIVE`, `INACTIVE`, `DUPLICATE`. Default: `ACTIVE`. Example: `--filter ACTIVE,DUPLICATE`. |
| `-q` | bool | no | Quiet mode: suppress all log output. Default: `false`. |
| `-o` | string | no | Write findings to the specified file. Banner and logs are not included; only the scan findings are written. |
| `--format` | string | no | Output format used when `-o` is set. Accepted values: `json` (default), `csv`, `raw` (plain table). |
| `<config-file>` | path | yes | Path to the TOML configuration file. |

### Execution Order

```
Parse flags → Load config.toml → Initialize scanners
  → Phase 1: Run prerequisite scanners concurrently (Syft SBOM generation)
  → Phase 2: Run dependent/independent scanners concurrently (Grype, KICS, OpenGrep)
           Any scanner whose prerequisite failed is skipped automatically
  → Load findings → [Sync to DefectDojo] → Display findings (stdout)
  → [-o: Dump findings to file in --format (json/csv/raw)]
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

# Directory to scan, relative to the SCAN_DIR environment variable.
# Used by both KICS and OpenGrep as the root of their scan targets.
path = "./my-service"

# KICS – infrastructure-as-code scanner
[kics]
# Infrastructure platform type. Passed as --type to KICS.
# Examples: "Dockerfile", "Terraform", "CloudFormation", "Kubernetes", "Ansible"
platform = "Dockerfile"
# Optional list of KICS query IDs to exclude from scanning.
# exclude_queries = ["a227ec01-f97a-4084-91a4-47b350c1db54"]

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
# Optional list of glob patterns passed to Syft via --exclude during SBOM generation.
# Use this to skip paths such as test sources from the SBOM (e.g. src/test/).
# Note: test-scoped pom.xml dependencies are not affected (Syft has no Maven scope filter).
# syft_exclude = ["**/src/test/**"]
# Depth to recursively resolve parent POMs (env: SYFT_JAVA_MAX_PARENT_RECURSIVE_DEPTH).
# 0 means no limit (the default).
# syft_max_parent_recursive_depth = 0

# OpenGrep – static application security testing (SAST) scanner.
[opengrep]
# Optional list of path patterns to exclude from scanning.
# exclude = ["**/vendor/**", "**/testdata/**"]
# Optional list of rule IDs to skip.
# exclude_rule = ["python.lang.security.audit.formatted-sql-query.formatted-sql-query"]

# Proxy – optional HTTP/HTTPS proxy settings forwarded to all scanner sub-processes.
# All fields are optional. Omit the entire section or leave fields empty to disable.
# [proxy]
# http_proxy    = "http://proxy.company.com:3128"
# https_proxy   = "http://proxy.company.com:3128"
# no_proxy      = "localhost,127.0.0.1"
# ssl_cert_file = "/path/to/ca.pem"   # PEM-encoded CA certificate (e.g. Burp Suite CA)
```

### Fields Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | no | Human-readable label; not used programmatically. |
| `protected_branches` | string array | no | Branches whose engagements get a 1-year end date. Defaults to empty (all branches get 1 week). |
| `path` | string | yes* | Path to the directory to scan. Resolved as `$SCAN_DIR/<path>`. Used by both KICS and OpenGrep. |
| `[kics].platform` | string | no | KICS platform filter (e.g. `Dockerfile`). When omitted KICS scans all supported types. |
| `[kics].exclude_queries` | string array | no | KICS query IDs to skip (e.g. `["a227ec01-f97a-4084-91a4-47b350c1db54"]`). |
| `[grype].ignore_states` | string | no | Comma-separated Grype vulnerability states to suppress (e.g. `not-fixed,unknown,wont-fix`). |
| `[grype].transitive_libraries` | bool | no | When `true`, Syft resolves transitive Java dependencies via Maven Central. Default: `false`. |
| `[grype].exclude` | string array | no | Path glob patterns to exclude from Grype scanning (e.g. `["**/vendor/**"]`). |
| `[grype].syft_exclude` | string array | no | Path glob patterns passed to Syft via `--exclude` during SBOM generation (e.g. `["**/src/test/**"]`). Excludes filesystem paths only; test-scoped `pom.xml` dependencies are unaffected. |
| `[grype].syft_max_parent_recursive_depth` | int | no | Maximum number of parent POM levels Syft will recursively resolve during Java/Maven analysis. `0` (default) means no limit. Forwarded as `SYFT_JAVA_MAX_PARENT_RECURSIVE_DEPTH`. |
| `[opengrep].exclude` | string array | no | Path glob patterns to exclude from OpenGrep scanning (e.g. `["**/vendor/**"]`). |
| `[opengrep].exclude_rule` | string array | no | OpenGrep rule IDs to skip (e.g. `["python.lang.security.audit.formatted-sql-query.formatted-sql-query"]`). |
| `[proxy].http_proxy` | string | no | HTTP proxy URL forwarded as `HTTP_PROXY` / `http_proxy` to all scanner sub-processes. |
| `[proxy].https_proxy` | string | no | HTTPS proxy URL forwarded as `HTTPS_PROXY` / `https_proxy` to all scanner sub-processes. |
| `[proxy].no_proxy` | string | no | Comma-separated list of hosts that bypass the proxy, forwarded as `NO_PROXY` / `no_proxy`. |
| `[proxy].ssl_cert_file` | string | no | Path to a PEM-encoded CA certificate bundle forwarded as `SSL_CERT_FILE` (Go tools) and `REQUESTS_CA_BUNDLE` (Python tools such as OpenGrep) to all scanner sub-processes. Required when using an intercepting proxy (e.g. Burp Suite). |

\* Required only if you want KICS or OpenGrep scanning to run. Omitting `path` while either `[kics]` or `[opengrep]` is configured will cause those scanners to use an empty path.

Omitting the entire `[grype]` section disables both Grype and the Syft SBOM generation step.

Omitting the entire `[opengrep]` section disables the SAST scanner.

Omitting the entire `[proxy]` section (or leaving all fields empty) disables proxy forwarding — scanner sub-processes inherit no proxy environment variables from this configuration.

Both the uppercase (`HTTP_PROXY`) and lowercase (`http_proxy`) variants of each proxy variable are set for maximum compatibility across tools. The `ssl_cert_file` value is emitted as both `SSL_CERT_FILE` (used by Go-based tools) and `REQUESTS_CA_BUNDLE` (used by Python-based tools such as OpenGrep).

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

ScopeGuardian uses a single DefectDojo **engagement** per project/branch combination to store all findings for that branch. Engagements are managed automatically — you never have to create or update them manually.

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

### Import vs Reimport

On the **first run** for a given engagement, each scanner calls the `/api/v2/import-scan/` endpoint. DefectDojo creates a new test and populates it with findings.

On **every subsequent run**, each scanner checks whether a test of the same scan type already exists in the engagement (via `GET /api/v2/tests/`). If one is found, it calls `/api/v2/reimport-scan/` instead. DefectDojo then:

- **Closes** findings that were present in the previous scan but are absent from the new one (`close_old_findings_product_scope=true`).
- **Creates** genuinely new findings as Active.
- **Does not reactivate** findings that a human previously suppressed (inactive, false-positive, accepted risk), even if they still appear in the scan (`do_not_reactivate=true`).

This means intentional triage decisions made in DefectDojo are preserved across pipeline runs.

### Finding Statuses

After uploading, ScopeGuardian fetches all findings for the engagement from DefectDojo and assigns each one a local status:

| Status | Condition | Meaning |
|--------|-----------|---------|
| `ACTIVE` | `active=true`, `duplicate=false` | Open, confirmed finding |
| `DUPLICATE` | `duplicate=true` | DefectDojo's deduplication engine identified an earlier occurrence |
| `INACTIVE` | `active=false` | Suppressed by a human (false-positive, accepted risk, manually closed) |

The `--filter` flag controls which statuses appear in the CLI output and in any file written with `-o`. By default only `ACTIVE` findings are shown. You can expand this with e.g. `--filter ACTIVE,DUPLICATE`.

The security gate (`--threshold`) always evaluates **only `ACTIVE` findings**, regardless of the `--filter` value.

### KICS Sync Behaviour

The KICS scanner uploads its JSON output file to DefectDojo as a `multipart/form-data` request. The endpoint used is `/api/v2/import-scan/` on the first run and `/api/v2/reimport-scan/` on subsequent runs. The following options are set on every upload:

| Option | Value | Effect |
|--------|-------|--------|
| Scan type | `KICS Scan` | Tells DefectDojo which parser to use |
| Severity threshold | `Info` | Import findings of all severities |
| Group by | `finding_title` | Merge findings with the same title |
| Create finding groups | `true` | Group related findings together |
| Apply tags to findings | `true` | Tag each finding with `IACST` |
| Close old findings | `true` | Findings absent from the new scan are closed automatically |
| Do not reactivate | `true` | Previously suppressed findings are not reactivated on reimport |
| Branch tag | `<branch>` | Associates the results with the scanned branch |

### Grype Sync Behaviour

The Grype scanner uploads its JSON output file to DefectDojo as a `multipart/form-data` request. The endpoint used is `/api/v2/import-scan/` on the first run and `/api/v2/reimport-scan/` on subsequent runs. The following options are set on every upload:

| Option | Value | Effect |
|--------|-------|--------|
| Scan type | `Anchore Grype` | Tells DefectDojo which parser to use |
| Severity threshold | `Info` | Import findings of all severities |
| Group by | `finding_title` | Merge findings with the same title |
| Create finding groups | `true` | Group related findings together |
| Apply tags to findings | `true` | Tag each finding with `SCA` |
| Close old findings | `true` | Findings absent from the new scan are closed automatically |
| Do not reactivate | `true` | Previously suppressed findings are not reactivated on reimport |
| Branch tag | `<branch>` | Associates the results with the scanned branch |

The **CWE/CVE** column in the CLI output and in file exports is populated with the CVE identifier (e.g. `CVE-2021-44228`) reported by Grype for each vulnerability match. KICS and OpenGrep findings use a CWE number in the same column.

### OpenGrep Sync Behaviour

The OpenGrep scanner uploads its JSON output file to DefectDojo as a `multipart/form-data` request. Before uploading, the file is enriched so that each result contains an `extra.severity` field required by DefectDojo's Semgrep JSON Report parser (the value is copied from `extra.metadata.impact`). The endpoint used is `/api/v2/import-scan/` on the first run and `/api/v2/reimport-scan/` on subsequent runs. The following options are set on every upload:

| Option | Value | Effect |
|--------|-------|--------|
| Scan type | `Semgrep JSON Report` | Tells DefectDojo which parser to use |
| Severity threshold | `Info` | Import findings of all severities |
| Group by | `finding_title` | Merge findings with the same title |
| Create finding groups | `true` | Group related findings together |
| Apply tags to findings | `true` | Tag each finding with `SAST` |
| Close old findings | `true` | Findings absent from the new scan are closed automatically |
| Do not reactivate | `true` | Previously suppressed findings are not reactivated on reimport |
| Branch tag | `<branch>` | Associates the results with the scanned branch |

### Security Gate with Sync

When both `--sync` and `--threshold` are set, the gate is evaluated against **DefectDojo's deduplicated findings** rather than the raw local scan output. This means:

- `INACTIVE` findings (suppressed, false-positive, accepted risk) are excluded from the count regardless of whether they appear in the latest scan.
- `DUPLICATE` findings do not inflate the count.
- Only `ACTIVE` findings surviving DefectDojo's deduplication and suppression logic are counted.

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

| Flags used | Findings evaluated by gate | Findings displayed / written |
|-----------|---------------------------|------------------------------|
| `--threshold` only | `ACTIVE` findings from local scan | controlled by `--filter` (default: `ACTIVE`) |
| `--threshold` + `--sync` | `ACTIVE` findings fetched from DefectDojo | controlled by `--filter` (default: `ACTIVE`) |

---

## Running with Docker

The provided `Dockerfile` builds a multi-stage image that bundles ScopeGuardian together with KICS, OpenGrep, Grype, and Syft.

```bash
# Build the image
docker build -t ScopeGuardian .

# Run a scan
docker run --rm \
  -v /path/to/your/project:/tmp/data/project \
  -v /path/to/config.toml:/config.toml \
  -e SCAN_DIR=/tmp/data \
  -e DD_URL=http://host.docker.internal:8080 \
  -e DD_ACCESS_TOKEN=<your-token> \
  ScopeGuardian \
  --projectName my-service \
  --branch main \
  --sync \
  /config.toml
```

Inside the container `SCAN_DIR` defaults to `/tmp/data`.

> **HTTPS proxy requirement — `--cap-add SYS_PTRACE`**
>
> OpenGrep bundles a Python interpreter inside its binary and reads `/proc/1/map_files` at startup to bootstrap it. In a Docker container this path requires the `CAP_SYS_PTRACE` Linux capability, which is dropped by default. If you run ScopeGuardian behind an HTTPS proxy (i.e. `[proxy].https_proxy` is set), you must add this capability so that OpenGrep can start:
>
> ```bash
> docker run --rm --cap-add SYS_PTRACE \
>   -v /path/to/your/project:/tmp/data/project \
>   -v /path/to/config.toml:/config.toml \
>   -e SCAN_DIR=/tmp/data \
>   -e DD_URL=http://host.docker.internal:8080 \
>   -e DD_ACCESS_TOKEN=<your-token> \
>   ScopeGuardian \
>   --projectName my-service \
>   --branch main \
>   --sync \
>   /config.toml
> ```
>
> With Docker Compose add `cap_add: [SYS_PTRACE]` to the ScopeGuardian service.
> With GitHub Actions use `options: --cap-add SYS_PTRACE` inside the `container:` block.

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

Before running ScopeGuardian with `--sync`, create a **Product** in DefectDojo whose name matches the `--projectName` value you will pass on the command line.

1. In DefectDojo go to **Products → Add Product**.
2. Set **Name** to exactly the value you will pass as `--projectName` (e.g. `my-service`).
3. Save. ScopeGuardian will manage engagements under this product automatically.
