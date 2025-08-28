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

func TestDumpHex(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "hexfile")
	content := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	os.WriteFile(tmp, content, 0644)
	defer os.Remove(tmp)

	var sb strings.Builder
	err := dumpHex(&sb, tmp)
	if err != nil {
		t.Fatal(err)
	}

	// The hex.Dump function adds a newline, which we want to match.
	expected := "00000000  01 02 03 04 05 06 07 08  09 0a 0b 0c 0d 0e 0f 10  |................|\n"
	if sb.String() != expected {
		t.Errorf("expected hex dump\n%q\ngot\n%q", expected, sb.String())
	}
}

func TestPrintTree(t *testing.T) {
	tmpdir := t.TempDir()
	// Create structure:
	// tmpdir/
	// ├── a.txt
	// └── b_dir/
	//     └── c.txt
	// Sorting of ReadDir is not guaranteed, but alphabetical is common.
	// Using names that are unlikely to be mis-ordered.
	os.WriteFile(filepath.Join(tmpdir, "a.txt"), []byte("a"), 0644)
	os.Mkdir(filepath.Join(tmpdir, "b_dir"), 0755)
	os.WriteFile(filepath.Join(tmpdir, "b_dir", "c.txt"), []byte("c"), 0644)

	var sb strings.Builder
	entries, err := os.ReadDir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}

	var dirs, files int
	// This test relies on a specific, common directory entry order.
	for i, entry := range entries {
		isLast := i == len(entries)-1
		printTree(&sb, filepath.Join(tmpdir, entry.Name()), "", isLast, "", true, &dirs, &files)
	}

	// Note: This expected output assumes 'a.txt' comes before 'b_dir'.
	expected := `├── a.txt
└── b_dir
    └── c.txt
`
	if sb.String() != expected {
		t.Errorf("expected tree:\n%q\ngot:\n%q", expected, sb.String())
	}
	if dirs != 1 {
		t.Errorf("expected 1 directory, got %d", dirs)
	}
	if files != 2 {
		t.Errorf("expected 2 files, got %d", files)
	}
}

func TestAllFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpdir := t.TempDir()
	outFile := filepath.Join(tmpdir, "out.md")

	// Create a file that should be ignored by default
	os.MkdirAll(filepath.Join(tmpdir, ".git"), 0755)
	os.WriteFile(filepath.Join(tmpdir, ".git", "config"), []byte("git stuff"), 0644)
	os.WriteFile(filepath.Join(tmpdir, "main.go"), []byte("go stuff"), 0644)

	// Case 1: Run WITHOUT --all
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"concat", "-o", outFile, tmpdir}
	main()

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if strings.Contains(content, "File Path: .git/config") {
		t.Errorf("--all=false: expected .git/config to be ignored, but it was included")
	}
	if !strings.Contains(content, "File Path: main.go") {
		t.Errorf("--all=false: expected main.go to be included, but it was not")
	}

	// Case 2: Run WITH --all
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"concat", "--all", "-o", outFile, tmpdir}
	main()

	data, err = os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	content = string(data)
	if !strings.Contains(content, "File Path: .git/config") {
		t.Errorf("--all=true: expected .git/config to be included, but it was ignored")
	}
}

func TestIncludeBinariesFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpdir := t.TempDir()
	outFile := filepath.Join(tmpdir, "out.md")

	// Create a text file and a binary file
	binFilePath := filepath.Join(tmpdir, "binary.bin")
	os.WriteFile(filepath.Join(tmpdir, "text.txt"), []byte("hello"), 0644)
	os.WriteFile(binFilePath, []byte{0x00, 0x01, 0x02}, 0644)

	// Case 1: Run WITHOUT --include-binaries
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"concat", "-o", outFile, tmpdir}
	main()

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	relPath, _ := filepath.Rel(tmpdir, binFilePath)
	if !strings.Contains(content, "File Path: "+relPath) {
		t.Errorf("expected binary file path to be present in output, but it wasn't")
	}
	if !strings.Contains(content, "[skipped binary file]") {
		t.Errorf("expected binary file to be skipped without the flag, but it wasn't")
	}
	if strings.Contains(content, "00000000") { // part of hex dump
		t.Errorf("expected binary file to NOT be dumped as hex, but it was")
	}

	// Case 2: Run WITH --include-binaries
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"concat", "--include-binaries", "-o", outFile, tmpdir}
	main()

	data, err = os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	content = string(data)
	if strings.Contains(content, "[skipped binary file]") {
		t.Errorf("expected binary file to be included with the flag, but it was skipped")
	}
	if !strings.Contains(content, "00000000") { // part of hex dump
		t.Errorf("expected binary file to be dumped as hex, but it wasn't")
	}
}
