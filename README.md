---

### 📘 New `README.md`

````markdown
# concat

`concat` is a CLI tool to flatten a project directory into a single Markdown file.  
This is useful for sharing multi-file projects with web-based AI assistants like ChatGPT.

---

## Features
- Generates a tree view of your directory
- Concatenates all file contents into one Markdown file
- Skips common junk files (`.git`, `.vscode`) by default
- Optionally include binary files as hex dumps
- Verbose mode to show progress
- Configurable output filename

---

## Install
```bash
git clone https://github.com/yourname/concat.git
cd concat
make build
````

This will place the binary at `bin/concat`.

---

## Usage

```bash
concat [options] <path>
```

### Options

* `-o, --output <file>` : Specify output file (default: `{parentDir}_OVERVIEW.md`)
* `--include-binaries` : Include binary files as hex dumps
* `--all`              : Include all files (ignore nothing)
* `--verbose`          : Show progress
* `--version`          : Show version and exit
* `-h, --help`         : Show help

---

## Examples

Concatenate current project:

```bash
concat .
```

Concatenate and include binary files:

```bash
concat --include-binaries .
```

Save output to a custom file:

```bash
concat -o project_OVERVIEW.md .
```

Verbose mode for large repos:

```bash
concat --all --verbose .
```

---

## Development

Run tests:

```bash
make test
```

Clean build:

```bash
make clean
```

Docker build:

```bash
make docker-build
```

---
