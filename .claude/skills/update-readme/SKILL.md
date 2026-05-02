---
name: update-readme
description: >
  Reviews recent code changes and updates README.md to keep it accurate.
  Invoke this skill after completing code changes that touch CLI flags,
  environment variables, configuration schema, scanner integrations,
  security-gate logic, DefectDojo sync behaviour, Docker image, or build
  requirements.
---

# Update README Skill

## When to invoke

Run this skill after any code change that touches one of the following areas:

| Area | Key files / packages |
|------|----------------------|
| CLI flags | `parser/` |
| Environment variables | `environnement_variable/` |
| Configuration file schema | `loader/`, `config.toml` |
| Scanner integration | `features/scans/*/`, `engine/` |
| Security gate logic | `features/security-gate/` |
| DefectDojo sync behaviour | `features/sync/`, `connectors/defectdojo/` |
| Docker image / binary paths | `Dockerfile` |
| Go module / build requirements | `go.mod`, `main.go` |

## Steps

1. **Read the changed files** — use the Read tool on every file you modified to recall exactly what changed.

2. **Check each README section** — open `README.md` and evaluate the sections below against your changes:

   | README section | What to verify |
   |----------------|----------------|
   | **Prerequisites** | Tool names, binary paths, and minimum Go version are accurate |
   | **Quick Start** | Example commands still work with the current flag names and env vars |
   | **CLI Usage** | Flag table is complete; names, types, and descriptions match `parser/` |
   | **Configuration File (`config.toml`)** | Every documented key exists; new keys are documented with examples |
   | **Environment Variables** | Table matches the vars loaded in `environnement_variable/` |
   | **How Engagements Are Handled** | Engagement naming and end-date rules are correct |
   | **How the Sync Feature Works** | Import vs reimport logic, flow description |
   | **How the Security Gate Works** | Threshold evaluation, severity ordering, exit code behaviour |
   | **Running with Docker** | Image name, tool paths, and `docker run` examples are current |

3. **Edit README.md** — for each section that is now inaccurate or missing information, apply the minimum edit needed to make it accurate. Do not rewrite sections that are still correct.

4. **Report what changed** — after editing, briefly summarise which README sections were updated and why (one line per section is enough).

5. **No-op case** — if none of the trigger areas were touched, state explicitly: "README.md does not require updates for this change."
