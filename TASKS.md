# Task List: concat Project Improvements

Each task should be implemented in its own **feature branch** and thoroughly tested before merging into `main`.

---

## 1. Version Management Improvement
**Branch:** `feature/versioning`

- [ ] **Remove hardcoded version from `main.go`**
  - Change:
    ```go
    var version = "v1.0.2-beta"
    ```
    → 
    ```go
    var version = "dev" // default if not injected
    ```
- [ ] **Update Makefile**
  - Add a `VERSION` variable:
    ```make
    VERSION ?= $(shell git describe --tags --always --dirty)
    build:
        go build -ldflags="-X main.version=$(VERSION)" -o bin/concat ./cmd/concat/main.go
    ```
- [ ] **Update release.yml workflow**
  - Inject version from Git tag:
    ```yaml
    VERSION=${GITHUB_REF#refs/tags/}
    go build -ldflags="-X main.version=$VERSION" -o dist/concat-${{ matrix.goos }}-${{ matrix.goarch }} ./cmd/concat
    ```

---

## 2. Go Version Compatibility
**Branch:** `feature/go-compatibility`

- [ ] Change `go.mod` from `go 1.24.6` → `go 1.22`.
- [ ] Update `go.yml` workflow matrix to test against Go 1.22 and 1.23.
- [ ] Run local build & test on Go 1.22 to confirm compatibility.

---

## 3. Path Normalization for Cross-Platform Consistency
**Branch:** `feature/path-normalization`

- [ ] Add helper function:
  ```go
  func toUnixPath(path string) string {
      return strings.ReplaceAll(path, string(os.PathSeparator), "/")
  }
- [ ] Wrap all file path outputs (`printTree`, concatenated file headers) with `toUnixPath`.
- [ ] Add tests to confirm `\` → `/` conversion works.

---

## 4. Expanded Test Coverage

**Branch:** `feature/test-coverage`

* [ ] **Tree output test**

  * Call `printTree` with a mock directory and compare output to expected tree.
* [ ] **CLI tests**

  * Manipulate `os.Args` to simulate CLI usage (`concat --help`, `concat --version`).
* [ ] **Error tests**

  * Non-existent directory → expect error.
  * Restricted permissions → expect handled error.
* [ ] **Path normalization tests**

  * Ensure `toUnixPath("a\\b\\c.txt") == "a/b/c.txt"`.

---

## 5. Enhanced Glob Pattern Matching

**Branch:** `feature/glob-matching`

* [ ] Update `shouldIncludeByGlob` to match **full relative path**, not just basename.
* [ ] Add support for recursive `**` patterns.

  * Option A: Use [`github.com/bmatcuk/doublestar`](https://pkg.go.dev/github.com/bmatcuk/doublestar).
  * Option B: Implement manually.
* [ ] Add tests to cover `**/*.go`, `docs/**`, etc.
* [ ] Update help text to clarify glob matching.

---

## 6. Default Ignores Visibility

**Branch:** `feature/ignored-visibility`

* [ ] Add `--show-ignored` flag (prints `defaultIgnores` and exits).
* [ ] Document default ignores in `--help`.
* [ ] (Optional) Add `--default-ignores "dir1,dir2,...` to override defaults.

---

## 7. Additional Considerations

**Branch:** `feature/extra-options`

* [ ] **Integration tests**

  * Run built binary against sample project → compare markdown output.
* [ ] **`--dry-run` flag**

  * Show what files *would* be processed without writing output.
* [ ] **Config file support**

  * Allow `.concat.yml` or `.concat.json` with ignore/include patterns.
* [ ] **File size limits**

  * Add `--max-file-size <bytes>` to skip large files.

---

## Implementation Order

1. `feature/versioning`
2. `feature/go-compatibility`
3. `feature/path-normalization`
4. `feature/test-coverage`
5. `feature/glob-matching`
6. `feature/ignored-visibility`
7. `feature/extra-options`

---

## Contribution Workflow

1. Create a new branch from `main`:

   ```bash
   git checkout main
   git pull origin main
   git checkout -b feature/<task-name>
   ```
2. Implement the changes and add tests.
3. Run:

   ```bash
   make build
   make test
   ```
4. Commit and push:

   ```bash
   git add .
   git commit -m "Implement <task-name>"
   git push origin feature/<task-name>
   ```
5. Open a Pull Request into `main`.

