# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- (Add new features here)

### Changed
- (Document changes in existing functionality here)

### Fixed
- (Document bug fixes here)

### Removed
- (Document removals/deprecations here)

---

## [v1.0.3] - 2025-08-25
### Added
- Version flag test (`--version`) to ensure version output is stable.
- CI builds now inject `ci-build` as version string.
- Release builds inject the Git tag as version string.

### Changed
- Version management is now handled via `-ldflags` instead of hardcoded variable.
- `Makefile` updated to allow `make build VERSION=x.y.z`.

---

## [v1.0.2-beta] - 2025-08-20
### Added
- Initial beta release.
- Directory tree rendering (`printTree`).
- Concatenation of project files into overview markdown.
- Glob include/exclude support.
- Ignore common junk files (`.git`, `.vscode`).

---

[Unreleased]: https://github.com/about80ninjas/concat/compare/v1.0.3...HEAD
[v1.0.3]: https://github.com/about80ninjas/concat/releases/tag/v1.0.3
[v1.0.2-beta]: https://github.com/about80ninjas/concat/releases/tag/v1.0.2-beta
