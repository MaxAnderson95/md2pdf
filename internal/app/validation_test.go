package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maxanderson95/md2pdf/internal/document"
)

func TestValidateLocalImagesRejectsMissingFile(t *testing.T) {
	t.Parallel()

	err := validateLocalImages("/tmp", []document.ImageRef{{Source: "missing.svg", Resolved: "file:///tmp/missing.svg"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidateLocalImagesAllowsExistingFileAndRemote(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	imagePath := filepath.Join(dir, "figure.svg")
	if err := os.WriteFile(imagePath, []byte("<svg></svg>"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	err := validateLocalImages(dir, []document.ImageRef{
		{Source: "figure.svg", Resolved: "file://" + filepath.ToSlash(imagePath)},
		{Source: "https://example.com/logo.png", Resolved: "https://example.com/logo.png", IsRemote: true},
	})
	if err != nil {
		t.Fatalf("validateLocalImages returned error: %v", err)
	}
}
