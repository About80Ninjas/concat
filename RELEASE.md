# Release Guide

This document explains how to cut a new release of **concat**.

---

## 🔹 1. Prerequisites

- You must have push access to the repo.
- Ensure your local `main` branch is up-to-date:
  ```bash
  git checkout main
  git pull origin main
---

## 🔹 2. Versioning Strategy

We use **Semantic Versioning (SemVer)**:

```
MAJOR.MINOR.PATCH[-prerelease]
```

* **MAJOR** → Breaking changes (incompatible CLI flags, output format changes).
* **MINOR** → Backward-compatible new features (new flags, options).
* **PATCH** → Backward-compatible bug fixes, small improvements, or documentation-only updates.
* **Prerelease** → Use `-alpha`, `-beta`, `-rc.1` for unstable builds.

Examples:

* `v1.0.3` → Stable patch release.
* `v1.1.0` → Adds new features.
* `v2.0.0` → Breaking change (e.g., CLI flags renamed).

---

## 🔹 3. Decide the Next Version

1. Check the current version:

   ```bash
   git describe --tags --abbrev=0
   ```

2. Based on the changes since last release, decide whether to bump:

   * **PATCH** (bugfixes, docs, small updates)
   * **MINOR** (new features, no breaking changes)
   * **MAJOR** (breaking changes)

---

## 🔹 4. Create the Release Tag

From `main`, tag the new version:

```bash
# Example: bumping to v1.0.3
git tag v1.0.3

# Push the tag to GitHub
git push origin v1.0.3
```

---

## 🔹 5. What Happens Next (CI/CD)

When a tag is pushed:

1. GitHub Actions (`release.yml`) will:

   * Build binaries for Linux, macOS, Windows (amd64 + arm64).
   * Inject the version string into the binary with `ldflags`.
   * Upload artifacts to the release.
2. A new GitHub Release will appear automatically.

---

## 🔹 6. Verify the Release

1. Go to the [Releases](../../releases) page.
2. Download a binary for your platform.
3. Check version:

   ```bash
   ./concat --version
   # Should print: concat version vX.Y.Z
   ```

---

## 🔹 7. Example Release Flow

* Finish **Task 2: Go version compatibility** → release `v1.0.4`.
* Add **Task 3: Path normalization** + **Task 5: Glob improvements** → release `v1.1.0`.
* Remove/rename CLI flags (breaking change) → release `v2.0.0`.

---

## 🔹 8. Pre-Releases

If you want to publish experimental versions:

```bash
git tag v1.1.0-beta.1
git push origin v1.1.0-beta.1
```

GitHub Actions will publish this as a **pre-release**.

---

## 🔹 9. Hotfixes

If you need to patch an older release:

1. Checkout the corresponding tag:

   ```bash
   git checkout v1.0.3
   ```
2. Create a hotfix branch:

   ```bash
   git checkout -b hotfix/v1.0.4
   ```
3. Apply fix, commit, merge, and tag `v1.0.4`.

---

## ✅ Summary

1. Update local `main`.
2. Decide the next version number.
3. Create and push a Git tag (`vX.Y.Z`).
4. GitHub Actions builds and publishes the release.
5. Verify the binaries and `--version` output.

---

