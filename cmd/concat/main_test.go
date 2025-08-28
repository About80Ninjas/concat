package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	if !shouldIgnore(".git") {
		t.Error("expected .git to be ignored")
	}
	if shouldIgnore("main.go") {
		t.Error("expected main.go not to be ignored")
	}
}

func TestIsBinary(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "binfile")
	os.WriteFile(tmp, []byte{0x00, 0xFF, 0x10}, 0644)
	defer os.Remove(tmp)

	isBin, err := isBinary(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if !isBin {
		t.Error("expected file to be detected as binary")
	}
}

func TestDumpText(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "textfile.txt")
	content := "hello world"
	os.WriteFile(tmp, []byte(content), 0644)
	defer os.Remove(tmp)

	var sb strings.Builder
	err := dumpText(&sb, tmp)
	if err != nil {
		t.Fatal(err)
	}
	if sb.String() != content {
		t.Errorf("expected %q got %q", content, sb.String())
	}
}

func TestIncludeExcludeGlobs(t *testing.T) {
	root := t.TempDir()

	// Create files
	paths := []string{
		filepath.Join(root, "keep.go"),
		filepath.Join(root, "keep.md"),
		filepath.Join(root, "skip.txt"),
		filepath.Join(root, "vendor", "junk.go"),
	}
	for _, p := range paths {
		os.MkdirAll(filepath.Dir(p), 0755)
		os.WriteFile(p, []byte("data"), 0644)
	}

	// include only *.go and *.md
	includeGlobs = []string{"*.go", "*.md"}
	excludeGlobs = []string{"vendor/*"}

	// should include keep.go
	if !shouldIncludeByGlob(root, paths[0]) {
		t.Errorf("expected %s to be included", paths[0])
	}
	// should include keep.md
	if !shouldIncludeByGlob(root, paths[1]) {
		t.Errorf("expected %s to be included", paths[1])
	}
	// should exclude skip.txt
	if shouldIncludeByGlob(root, paths[2]) {
		t.Errorf("expected %s to be excluded by include filter", paths[2])
	}
	// should exclude vendor/junk.go
	if !shouldExcludeByGlob(root, paths[3]) {
		t.Errorf("expected %s to be excluded by exclude filter", paths[3])
	}
}

func TestDefaultIncludeBehavior(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "anyfile.txt")
	os.WriteFile(path, []byte("data"), 0644)

	// no includeGlobs set -> everything allowed
	includeGlobs = []string{}
	excludeGlobs = []string{}

	if !shouldIncludeByGlob(root, path) {
		t.Errorf("expected %s to be included by default", path)
	}
}

func TestGoalFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	tmpdir := t.TempDir()
	outFile := filepath.Join(tmpdir, "out.md")

	os.Args = []string{"concat", "--goal", "Test project goal", "-o", outFile, tmpdir}

	main()

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "Project Summary & Goal") {
		t.Errorf("expected Project Summary & Goal section, got:\n%s", content)
	}
	if !strings.Contains(content, "Test project goal") {
		t.Errorf("expected goal text in output, got:\n%s", content)
	}
}

func TestWithContextFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	tmpdir := t.TempDir()
	outFile := filepath.Join(tmpdir, "out.md")

	// Run with --with-context (git/make may fail, but output should mention failure)
	os.Args = []string{"concat", "--with-context", "-o", outFile, tmpdir}
	main()

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "Project Context") {
		t.Errorf("expected Project Context section, got:\n%s", content)
	}
	// Ensure at least one subheader appears
	if !strings.Contains(content, "## Git Status") {
		t.Errorf("expected Git Status section, got:\n%s", content)
	}
}

func TestDetectLang(t *testing.T) {
	cases := map[string]string{
		"file.go":   "go",
		"config.yml": "yaml",
		"config.yaml": "yaml",
		"data.json": "json",
		"readme.md": "markdown",
		"script.sh": "bash",
		"script.ps1": "powershell",
		"config.toml": "toml",
		"unknown.xyz": "",
	}
	for file, want := range cases {
		got := detectLang(file)
		if got != want {
			t.Errorf("detectLang(%q) = %q, want %q", file, got, want)
		}
	}
}

func TestVersionFlag(t *testing.T) {
	// Save old args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Simulate running: concat --version
	os.Args = []string{"concat", "--version"}

	// Capture stdout
	var buf strings.Builder
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan struct{})
	go func() {
		io.Copy(&buf, r)
		close(done)
	}()

	// Run main()
	main()

	w.Close()
	os.Stdout = old
	<-done

	output := buf.String()

	// Updated regex pattern to match "concat version v0.0.0" format with optional pre-release suffixes
	versionPattern := regexp.MustCompile(`concat version (?:v\d+\.\d+\.\d+(?:-[a-zA-Z0-9]+)*|dev)`)

	if !versionPattern.MatchString(output) {
		t.Errorf("expected version output to match pattern 'concat version v0.0.0' (with optional pre-release suffix), got %q", output)
	}
}