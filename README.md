# concat

[![Go CI](https://github.com/about80ninjas/concat/actions/workflows/test.yml/badge.svg)](https://github.com/about80ninjas/concat/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/about80ninjas/concat)](https://goreportcard.com/report/github.com/about80ninjas/concat)
[![Go Reference](https://pkg.go.dev/badge/github.com/about80ninjas/concat.svg)](https://pkg.go.dev/github.com/about80ninjas/concat)


`concat` is a CLI tool to flatten a project directory into a single Markdown file.  
This is useful for sharing multi-file projects with web-based AI assistants like ChatGPT.

---

## Features
- Generates a tree view of your directory
- Concatenates all file contents into one Markdown file
- Skips common junk files (`.git`, `.vscode`) by default
- Glob filters (`--include`, `--exclude`) for fine control
- Optionally include binary files as hex dumps
- Verbose mode to show progress
- Configurable output filename

---

## Install
```bash
git clone https://github.com/your-username/concat.git
cd concat
go mod init github.com/your-username/concat
make build
