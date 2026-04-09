# Changelog

All notable changes to Scope-Guardian are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial open-source release of Scope-Guardian
- KICS infrastructure-as-code scanner integration
- Grype software-composition analysis (SCA) scanner integration with Syft SBOM generation
- OpenGrep SAST scanner integration
- DefectDojo sync via `--sync` flag (engagement auto-create/reuse, per-scanner import)
- Security gate via `--threshold` flag with configurable severity/count rules
- Quiet mode (`-q`) and file output (`-o`) flags
- Multi-stage Docker image bundling all scanner binaries
- `docker-compose.yml` for local DefectDojo setup

[Unreleased]: https://github.com/ParanoiHack/scope-guardian/compare/HEAD...HEAD
