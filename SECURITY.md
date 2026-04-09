# Security Policy

## Supported Versions

The following versions of ScopeGuardian currently receive security updates:

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |

## Reporting a Vulnerability

We take security seriously. If you discover a vulnerability in ScopeGuardian, **please do not open a public GitHub issue**.

Report it privately using one of these channels:

- **Email**: contact@paranoihack.com
- **GitHub private advisory**: [Open a private advisory](https://github.com/ParanoiHack/ScopeGuardian/security/advisories/new)

### What to include

- A description of the vulnerability and its potential impact
- Steps to reproduce (proof-of-concept commands, configuration, or code)
- Affected version(s) of ScopeGuardian
- Any suggested mitigations or fixes

### Response timeline

| Step          | Target time |
| ------------- | ----------- |
| Acknowledgement | 48 hours |
| Severity assessment | 7 days |
| Patch for critical/high issues | 30 days |
| Public disclosure | After patch is available |

## Scope

This policy covers the `ScopeGuardian` CLI source code in this repository. It does **not** cover third-party tools bundled in the Docker image (KICS, OpenGrep, Grype, Syft). Please report vulnerabilities in those tools to their respective maintainers.

## Credits

We will publicly credit reporters in the release notes unless anonymity is requested. Thank you for helping keep ScopeGuardian secure.
