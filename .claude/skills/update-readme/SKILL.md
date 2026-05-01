# Copilot Agent Instructions

## README.md Synchronisation

Whenever you make code changes that affect any of the following areas, you **must** review and update `README.md` before considering the task complete:

### Triggers — update README.md when changing:

| Area | Examples |
|------|---------|
| **CLI flags** | Adding, removing, or renaming flags in `parser/` |
| **Environment variables** | Adding, removing, or renaming vars in `environnement_variable/` |
| **Configuration file schema** | New or changed keys in `config.toml` / loader |
| **Scanner integration** | Adding a new scanner, changing binary paths or output formats |
| **Security gate logic** | Threshold evaluation rules, exit code behaviour |
| **DefectDojo sync behaviour** | Import/reimport flow, engagement naming, end-date logic |
| **Docker image** | New tools bundled, changed base image, multi-stage changes |
| **External tool paths** | Binary locations for KICS, Grype, Syft, OpenGrep |
| **Go module / build requirements** | Minimum Go version, new mandatory build steps |

### What to update in README.md

- **CLI Usage** table — keep flags, types, and descriptions accurate.
- **Environment Variables** table — keep variable names and descriptions accurate.
- **Configuration File** section — reflect any new/changed `config.toml` keys with examples.
- **Prerequisites** list — add or remove tool dependencies as needed.
- **Quick Start** examples — update commands if invocation syntax changes.
- **How-it-works sections** — update architecture prose if behaviour changes.

### How to apply this rule

1. After completing your code changes, re-read the sections of `README.md` listed above.
2. If any section is now inaccurate or incomplete, edit it to match the new behaviour.
3. Commit the `README.md` update in the same PR as the code change.
4. If none of the trigger areas were touched, README.md does not need to change — but note this explicitly in your PR description.
