
# concat

[![Go CI](https://github.com/about80ninjas/concat/actions/workflows/go.yml/badge.svg)](https://github.com/about80ninjas/concat/actions/workflows/go.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/about80ninjas/concat)](https://github.com/about80ninjas/concat/releases/latest)

`concat` is a CLI tool to flatten a project directory into a single Markdown file. This is useful for sharing multi-file projects with web-based AI assistants like ChatGPT.

---

## Features

-   Generates a tree view of your directory.
-   Concatenates all file contents into one Markdown file.
-   Skips common junk files (`.git`, `.vscode`) by default.
-   Glob include/exclude support (see Options for details on current behavior).
-   Optionally include **project summary/goal** at the top of the overview.
-   Optionally include **git + build context** to capture project state.
-   Syntax highlighting in file content sections based on file extension.
 -   Optionally include binary files as hex dumps.
 -   Verbose mode to show progress.
 -   Configurable output filename.

---

## Installation

### From Pre-compiled Binaries

You can download the latest pre-compiled binaries from the [releases page](https://github.com/about80ninjas/concat/releases/latest) or use the following `curl` commands for your platform.

**Linux (amd64):**

```bash
curl -L https://github.com/about80ninjas/concat/releases/latest/download/concat-linux-amd64 -o concat
chmod +x concat
sudo mv concat /usr/local/bin/
````

**Linux (arm64):**

```bash
curl -L https://github.com/about80ninjas/concat/releases/latest/download/concat-linux-arm64 -o concat
chmod +x concat
sudo mv concat /usr/local/bin/
```

**macOS (amd64 - Intel):**

```bash
curl -L https://github.com/about80ninjas/concat/releases/latest/download/concat-darwin-amd64 -o concat
chmod +x concat
sudo mv concat /usr/local/bin/
```

**macOS (arm64 - Apple Silicon):**

```bash
curl -L https://github.com/about80ninjas/concat/releases/latest/download/concat-darwin-arm64 -o concat
chmod +x concat
sudo mv concat /usr/local/bin/
```

**Windows (amd64):**
(Using PowerShell)

```powershell
curl -L https://github.com/about80ninjas/concat/releases/latest/download/concat-windows-amd64.exe -o concat.exe
# Move concat.exe to a directory in your PATH
```

**Windows (arm64):**
(Using PowerShell)

```powershell
curl -L https://github.com/about80ninjas/concat/releases/latest/download/concat-windows-arm64.exe -o concat.exe
# Move concat.exe to a directory in your PATH
```

### From Source

If you have Go installed, you can build `concat` from source.
Requires Go 1.24.6 or later.

```bash
git clone https://github.com/about80ninjas/concat.git
cd concat
make build
# The binary will be at ./bin/concat
```

-----

## Usage

### Basic Usage

To run `concat` on the current directory and generate an overview file:

```bash
concat .
```

This will create a file named `parentDir_OVERVIEW.md` in the current directory.

### Options

Here are the available command-line options:

| Flag                 | Shorthand | Description                                           | Default                             |
| -------------------- | --------- | ----------------------------------------------------- | ----------------------------------- |
| `--output <file>`    | `-o`      | Specify the output file name.                         | `{parentDir}_OVERVIEW.md`           |
| `--include-binaries` |           | Include binary files as hex dumps.                    | `false`                             |
| `--all`              |           | Include all files, ignoring default ignore patterns.  | `false`                             |
| `--include <globs>`  |           | Comma-separated glob patterns to include (matches against file base name). | `""` (all files)                    |
| `--exclude <globs>`  |           | Comma-separated glob patterns to exclude (matches against relative path). | `""` (no exclusions)                |
| `--verbose`          |           | Show progress while scanning.                         | `false`                             |
| `--goal <text>`      |           | Add a "Project Summary & Goal" section at the top.    | `""` (disabled)                      |
| `--with-context`     |           | Include git status, recent commits, and build context.| `false`                             |
| `--version`          |           | Show the version and exit.                            |                                     |
| `--help`             | `-h`      | Show the help message.                                |                                     |

### Examples

**Specify an output file:**

```bash
concat -o project_overview.md .
```

**Include binary files:**

```bash
concat --include-binaries ./my_project
```

**Include only Go and Markdown files, excluding the `vendor` directory:**

```bash
concat --include "*.go,*.md" --exclude "vendor/*" .
```

**Include all files and show verbose output:**

```bash
concat --all --verbose --output overview.md .
```

**Add a project goal at the top of the overview:**

```bash
concat --goal "Refactor file path normalization" .
```

**Include git status, recent commits, and build/test context:**

```bash
concat --with-context .
```

**Both goal and context together:**

```bash
concat --goal "Improve test coverage" --with-context .
```

-----

## Contributing

Contributions are welcome\! Please feel free to submit a pull request.

-----

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
