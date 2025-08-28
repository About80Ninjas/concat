package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var version = "dev" // default, overridden by -ldflags

// default ignore patterns (dirs & files)
var defaultIgnores = []string{
	".git",
	".vscode",
}

var (
	verbose      bool
	includeGlobs []string
	excludeGlobs []string
	goal         string
	withContext  bool
)

func logVerbose(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf(msg+"\n", args...)
	}
}

func showHelp() {
	fmt.Print(`concat - Concatenate files in a directory tree into a single Markdown overview

Usage:
  concat [options] <path>

Options:
  -o, --output <file>       Specify output file (default: {parentDir}_OVERVIEW.md)
  --include-binaries        Include binary files as hex dumps (default: false)
  --all                     Include all files (ignore nothing)
  --include <globs>         Comma-separated glob patterns to include (e.g. "*.go,*.md")
  --exclude <globs>         Comma-separated glob patterns to exclude (e.g. "*.log,vendor/*")
  --verbose                 Show progress while scanning
  --version                 Show version and exit
  -h, --help                Show this help message

Description:
  This tool walks through the provided directory, builds a tree-like structure
  of all files and subdirectories, and writes a Markdown file.

Examples:
  concat .
  concat -o project_OVERVIEW.md .
  concat --include-binaries ./my_project
  concat --all --verbose --output overview.md .
  concat --include "*.go,*.md" --exclude "vendor/*" .
`)
}

func main() {
	var outFileFlag string
	var includeBinaries bool
	var includeAll bool
	var showVersion bool
	var includeStr, excludeStr string

	flag.StringVar(&outFileFlag, "o", "", "output file")
	flag.StringVar(&outFileFlag, "output", "", "output file")
	flag.BoolVar(&includeBinaries, "include-binaries", false, "include binary files as hex dumps")
	flag.BoolVar(&includeAll, "all", false, "include all files")
	flag.BoolVar(&verbose, "verbose", false, "show progress")
	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.StringVar(&includeStr, "include", "", "comma-separated glob patterns to include")
	flag.StringVar(&excludeStr, "exclude", "", "comma-separated glob patterns to exclude")
	flag.StringVar(&goal, "goal", "", "project goal to include in overview")
	flag.BoolVar(&withContext, "with-context", false, "include git/build context")
	flag.Usage = showHelp
	flag.Parse()

	if showVersion {
		fmt.Println("concat version", version)
		return
	}

	if flag.NArg() != 1 {
		showHelp()
		os.Exit(1)
	}

	if includeStr != "" {
		includeGlobs = strings.Split(includeStr, ",")
	}
	if excludeStr != "" {
		excludeGlobs = strings.Split(excludeStr, ",")
	}

	root := flag.Arg(0)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Println("Error resolving path:", err)
		os.Exit(1)
	}

	parentName := filepath.Base(absRoot)
	outputFile := outFileFlag
	if outputFile == "" {
		outputFile = filepath.Join(absRoot, fmt.Sprintf("%s_OVERVIEW.md", parentName))
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer out.Close()

	// Project Summary & Goal
	if goal != "" {
		fmt.Fprintln(out, "# Project Summary & Goal")
		fmt.Fprintf(out, "- **Goal:** %s", goal)
		fmt.Fprintln(out, "- **Project:** concat, a CLI tool to flatten a project directory into a single Markdown file for AI assistants.")
		fmt.Fprintln(out, "- **Language:** Go (version 1.24.6)")
		fmt.Fprintln(out, "---")
	}

	// Project Context
	if withContext {
		fmt.Fprintln(out, "# Project Context")
		runAndWrite(out, "## Git Status", "git", "status")
		runAndWrite(out, "## Recent Commits", "git", "log", "-n", "3", "--oneline")
		runAndWrite(out, "## Test Commands", "make", "test")
		runAndWrite(out, "## Build Commands", "make", "build")
		fmt.Fprintln(out, "---")
	}

	// 1. Directory Tree
	fmt.Fprintln(out, "# Directory Structure")
	dirs, files := 0, 0
	fmt.Fprintln(out, ".") // root marker

	entries, err := os.ReadDir(absRoot)
	if err != nil {
		fmt.Println("Error reading root:", err)
		os.Exit(1)
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1
		printTree(out, filepath.Join(absRoot, entry.Name()), "", isLast, outputFile, includeAll, &dirs, &files)
	}

	fmt.Fprintf(out, "\n%d directories, %d files\n", dirs, files)

	// 2. Concat file contents
	fmt.Fprintln(out, "# File Contents")
	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == outputFile {
			return nil
		}
		if !includeAll && shouldIgnore(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldExcludeByGlob(absRoot, path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldIncludeByGlob(absRoot, path) {
			if !info.IsDir() {
				relPath, _ := filepath.Rel(absRoot, path)
				logVerbose("Processing %s", relPath)

				isBin, err := isBinary(path)
				if err != nil {
					return nil // unreadable → skip
				}

				lang := detectLang(path)
				fmt.Fprintf(out, "\n-----\nFile Path: %s\n\n", relPath)

				if isBin {
					if includeBinaries {
						err = dumpHex(out, path)
						if err != nil {
							return err
						}
					} else {
						fmt.Fprintln(out, "[skipped binary file]")
					}
				} else {
					if lang != "" {
						fmt.Fprintf(out, "```%s\n", lang)
					}
					err = dumpText(out, path)
					if err != nil {
						return err
					}
					if lang != "" {
						fmt.Fprintln(out, "```")
					}
				}

				fmt.Fprintln(out, "-----")
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error concatenating files:", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s\n", outputFile)
}

// printTree prints a nice tree like `tree` command
func printTree(out io.Writer, path, prefix string, isLast bool, outputFile string, includeAll bool, dirCount, fileCount *int) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if path == outputFile {
		return nil
	}
	if !includeAll && shouldIgnore(info.Name()) {
		if info.IsDir() {
			return nil
		}
		return nil
	}
	if shouldExcludeByGlob(filepath.Dir(outputFile), path) {
		return nil
	}

	branch := "├── "
	if isLast {
		branch = "└── "
	}
	fmt.Fprintf(out, "%s%s%s\n", prefix, branch, info.Name())

	if info.IsDir() {
		*dirCount++
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for i, entry := range entries {
			isLastEntry := i == len(entries)-1
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(out, filepath.Join(path, entry.Name()), newPrefix, isLastEntry, outputFile, includeAll, dirCount, fileCount)
		}
	} else {
		*fileCount++
	}
	return nil
}

// shouldIgnore checks if file/dir should be ignored by default
func shouldIgnore(name string) bool {
	for _, pat := range defaultIgnores {
		if strings.EqualFold(name, pat) {
			return true
		}
	}
	return false
}

func shouldIncludeByGlob(root, path string) bool {
	if len(includeGlobs) == 0 {
		return true // default: include all
	}
	rel, _ := filepath.Rel(root, path)
	for _, pat := range includeGlobs {
		match, _ := filepath.Match(pat, filepath.Base(rel))
		if match {
			return true
		}
	}
	return false
}

func shouldExcludeByGlob(root, path string) bool {
	if len(excludeGlobs) == 0 {
		return false
	}
	rel, _ := filepath.Rel(root, path)
	for _, pat := range excludeGlobs {
		match, _ := filepath.Match(pat, rel)
		if match {
			return true
		}
	}
	return false
}

// isBinary checks if file likely contains binary data
func isBinary(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, 8000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	buf = buf[:n]

	if bytes.Contains(buf, []byte{0}) {
		return true, nil
	}
	if !utf8.Valid(buf) {
		return true, nil
	}
	return false, nil
}

// dumpText copies text file contents into output
func dumpText(out io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(out, f)
	return err
}

// dumpHex writes a hex dump of a binary file
func dumpHex(out io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 16)
	writer := bufio.NewWriter(out)

	for {
		n, err := f.Read(buf)
		if n > 0 {
			fmt.Fprint(writer, hex.Dump(buf[:n]))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// runAndWrite executes a shell command and writes output with a header
func runAndWrite(out io.Writer, header string, cmd string, args ...string) {
	fmt.Fprintln(out, header)
	c := exec.Command(cmd, args...)
	b, err := c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(out, "(failed to run %s: %v)\n", cmd, err)
	} else {
		fmt.Fprintln(out, string(b))
	}
	fmt.Fprintln(out)
}

// detectLang returns a markdown code fence language based on file extension
func detectLang(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "go"
	case ".yml", ".yaml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	case ".sh":
		return "bash"
	case ".ps1":
		return "powershell"
	case ".toml":
		return "toml"
	default:
		return ""
	}
}
