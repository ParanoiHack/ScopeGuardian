---
name: go-fmt
description: >
  Runs go fmt on all Go source files that were modified as part of the current
  code change. Invoke this skill after completing any code changes to Go files
  to ensure consistent formatting before committing.
---

# Go Fmt Skill

## When to invoke

Run this skill after any code change that creates or modifies one or more `.go` files.

## Steps

1. **Identify changed Go files** — run the following command to list every `.go` file that differs from the last commit:

   ```bash
   git diff --name-only HEAD | grep '\.go$'
   ```

   If there are no staged/unstaged changes yet (e.g. files were just written), list all `.go` files that were touched during this session instead.

2. **Run go fmt** — format all changed files in one call:

   ```bash
   go fmt ./...
   ```

   This is preferred over per-file invocation because it also catches any indirectly affected files in the same packages.

3. **Verify the result** — confirm that `go fmt` exited with code `0`. If it exited with a non-zero code, read the error output, fix the root cause (usually a syntax error in the modified file), and re-run `go fmt ./...`.

4. **Check for reformatted files** — run:

   ```bash
   git diff --name-only
   ```

   Any files listed here were reformatted by `go fmt`. This is expected and correct.

5. **Report the outcome** — summarise what happened in one or two lines:
   - If files were reformatted: list them and note that formatting was applied.
   - If no files changed after `go fmt`: state "All changed Go files were already correctly formatted."
   - If an error occurred: report the error message and the file it came from.

## Notes

- `go fmt` is a no-op on already-correctly-formatted code, so it is always safe to run.
- Do **not** skip this skill for test files (`*_test.go`) — they must be formatted too.
- The module root for this repository is `/home/runner/work/ScopeGuardian/ScopeGuardian` (or the directory where `go.mod` lives). Always run `go fmt` from that directory so the `./...` pattern covers all packages.
