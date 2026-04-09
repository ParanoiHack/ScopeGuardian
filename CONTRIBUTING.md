# Contributing to Scope-Guardian

Thank you for taking the time to contribute! Every contribution — bug reports, feature ideas, documentation improvements, and code changes — is welcome and appreciated.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [How to Contribute](#how-to-contribute)
4. [Development Setup](#development-setup)
5. [Coding Standards](#coding-standards)
6. [Commit Message Guidelines](#commit-message-guidelines)
7. [License](#license)

---

## Code of Conduct

By participating in this project you agree to abide by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

---

## Getting Started

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally:

   ```bash
   git clone https://github.com/<your-username>/scope-guardian.git
   cd scope-guardian
   ```

3. **Add the upstream remote** so you can keep your fork in sync:

   ```bash
   git remote add upstream https://github.com/ParanoiHack/scope-guardian.git
   ```

---

## How to Contribute

### Reporting Bugs

- Search the [existing issues](https://github.com/ParanoiHack/scope-guardian/issues) first to avoid duplicates.
- If no issue exists, open one using the **Bug Report** template and fill in all fields.
- For **security vulnerabilities** follow the [Security Policy](SECURITY.md) instead of opening a public issue.

### Requesting Features

- Search [existing issues](https://github.com/ParanoiHack/scope-guardian/issues) for similar ideas first.
- If none exists, open an issue using the **Feature Request** template.

### Submitting Pull Requests

1. Create a branch from `main`:

   ```bash
   git checkout -b feat/my-feature
   ```

2. Make your changes (see [Development Setup](#development-setup) below).
3. Run tests and the linter before pushing:

   ```bash
   go test ./...
   go vet ./...
   ```

4. Commit your changes following the [commit guidelines](#commit-message-guidelines).
5. Push your branch and open a pull request against `main`, filling in the pull request template completely.

A maintainer will review your PR. Please be patient — reviews may take a few days.

---

## Development Setup

**Prerequisites**

- Go 1.24+
- Docker (optional — for building and testing the full Docker image)

**Build from source**

```bash
go build -o scope-guardian .
```

**Run the test suite**

```bash
go test ./...
```

**Copy the environment template**

```bash
cp .env.example .env
# Fill in DD_URL, DD_ACCESS_TOKEN, and other variables as needed
```

**Start a local DefectDojo instance for integration testing**

```bash
docker compose up -d
```

---

## Coding Standards

- Follow [Effective Go](https://go.dev/doc/effective_go) conventions.
- Keep functions focused; prefer small, composable units.
- All exported symbols must have a Go doc comment.
- New scanner integrations should follow the existing pattern under `features/scans/`.
- Never commit secrets, credentials, or `.env` files.

---

## Commit Message Guidelines

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short summary>
```

Common types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `ci`.

**Examples**

```
feat(grype): add CVSS score filtering support
fix(engine): skip scanner when prerequisite fails
docs: update README threshold examples
chore: bump go-pretty to v6.7.7
```

---

## License

By contributing you agree that your contributions will be licensed under the [MIT License](LICENSE).
