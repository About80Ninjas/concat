# Task List: concat Project Improvements

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


