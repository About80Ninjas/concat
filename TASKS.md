# Task List: concat Project Improvements

## ! FIXES
- [ ] Add a "Project Summary & Goal" Section at the Top - Example: --goal "My goal is to..."
```markdown
# Project Summary & Goal

- **Goal:** Refactor the file path normalization to ensure cross-platform consistency.
- **Project:** `concat`, a CLI tool to flatten a project directory into a single Markdown file for AI assistants.
- **Language:** Go (version 1.24.6)
- **Key Files:**
  - `cmd/concat/main.go`: Main application logic.
  - `Makefile`: Build and test commands.
  - `.github/workflows/release.yml`: Release automation.

---
```
- [ ] Include High-Level Git & Build Context - Your concat tool can execute git status, git log -n 3 --oneline, and go list -m all and pipe the output into this section.
```markdown
# Project Context

## Git Status
On branch `feat/path-normalization`
Your branch is up to date with 'origin/feat/path-normalization'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   cmd/concat/main.go

no changes added to commit (use "git add" and/or "git commit -a")

## Recent Commits
- `a1b2c3d (HEAD -> feat/path-normalization)` Add initial path helper function
- `d4e5f6a` Refactor glob matching logic
- `g7h8i9j (origin/main, main)` Update release documentation

## Build & Test Commands
- `make build`: Builds the binary to `./bin/concat`.
- `make test`: Runs all tests.

---
```
- [ ] Enhance File Content Formatting -When writing file contents, detect the file type from its extension (.go, .yml, .md) and add it to the Markdown code block fence (e.g., ```go).
```markdown
# File Contents

## Core Application Logic

-----
File Path: cmd/concat/main.go

` ` `go
package main

import (
    "fmt"
    // ...
)
// ...
` ` `
-----

## CI/CD Workflows

-----
File Path: .github/workflows/go.yml

` ` `yaml
permissions:
  contents: read
name: Go CI
# ...
` ` `
-----
File Path: .github/workflows/release.yml

` ` `yaml
name: Release
# ...
` ` `
-----
```

## 1. Version Management Improvement ✅ Completed
- [x] Remove hardcoded version from `main.go` (default → `"dev"`)
- [x] Update Makefile to support `VERSION` variable with `-ldflags`
- [x] Update `release.yml` to inject version from Git tag
- [x] Update `go.yml` to inject `ci-build` during CI
- [x] Add test for `--version` flag (basic output check)

---

## 2. Go Version Compatibility (Next Task)
- [ ] Change `go.mod` from `go 1.24.6` → `go 1.22`
- [ ] Update `go.yml` workflow matrix to test against Go 1.22 and 1.23
- [ ] Run local build & test on Go 1.22 to confirm compatibility

---

## 3. Path Normalization for Cross-Platform Consistency
- [ ] Add helper function:
  ```go
  func toUnixPath(path string) string {
      return strings.ReplaceAll(path, string(os.PathSeparator), "/")
  }
* [ ] Wrap all file path outputs (`printTree`, concatenated file headers) with `toUnixPath`
* [ ] Add tests for Windows-style path conversion

---

## 4. Expanded Test Coverage

* [ ] Tree output test
* [ ] CLI tests (`os.Args`, including `--help`, `--version`)
* [ ] Error tests (invalid paths, permissions)
* [ ] Path normalization tests

---

## 5. Enhanced Glob Pattern Matching

* [ ] Update `shouldIncludeByGlob` to match **full relative path**
* [ ] Add support for `**` recursive globs
* [ ] Add tests for `**/*.go`, `docs/**`
* [ ] Update help text with glob details

---

## 6. Default Ignores Visibility

* [ ] Add `--show-ignored` flag
* [ ] Document default ignores in help output
* [ ] (Optional) Add `--default-ignores` override

---

## 7. Additional Considerations

* [ ] Integration tests (binary end-to-end)
* [ ] `--dry-run` flag
* [ ] Config file support
* [ ] File size limit option
* [ ] Ignore the concat binary by default

---

## Implementation Order

1. ✅ Version Management Improvement
2. Go Version Compatibility
3. Path Normalization
4. Expanded Test Coverage
5. Enhanced Glob Pattern Matching
6. Default Ignores Visibility
7. Additional Features

---


