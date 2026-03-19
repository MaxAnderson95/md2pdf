package input

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveRejectsNonMarkdown(t *testing.T) {
	t.Parallel()

	_, err := Resolve("notes.txt")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveFindsNextAvailableOutput(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputPath := filepath.Join(dir, "guide.md")
	if err := os.WriteFile(inputPath, []byte("# guide\n"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "guide.pdf"), []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "guide (01).pdf"), []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	paths, err := Resolve(inputPath)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	expected := filepath.Join(dir, "guide (02).pdf")
	if paths.OutputPath != expected {
		t.Fatalf("unexpected output path: got %q want %q", paths.OutputPath, expected)
	}
}

func TestNextOutputPathReturnsPrimaryWhenAvailable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	outputPath, err := NextOutputPath(dir, "report")
	if err != nil {
		t.Fatalf("NextOutputPath returned error: %v", err)
	}
	expected := filepath.Join(dir, "report.pdf")
	if outputPath != expected {
		t.Fatalf("unexpected output path: got %q want %q", outputPath, expected)
	}
}
