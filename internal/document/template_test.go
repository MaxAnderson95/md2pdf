package document

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildIncludesBaseHrefAndImageRefs(t *testing.T) {
	t.Parallel()

	builder, err := NewBuilder()
	if err != nil {
		t.Fatalf("NewBuilder returned error: %v", err)
	}

	inputPath := filepath.Join("/tmp", "guide.md")
	doc, err := builder.Build("<h1>Guide</h1><p><img src=\"images/diagram.svg\"></p>", inputPath)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !strings.Contains(doc.HTML, "<base href=\"file:///tmp/\">") {
		t.Fatalf("document HTML missing base href: %s", doc.HTML)
	}
	if len(doc.ImageRefs) != 1 {
		t.Fatalf("unexpected image refs: %+v", doc.ImageRefs)
	}
	if doc.ImageRefs[0].Resolved != "file:///tmp/images/diagram.svg" {
		t.Fatalf("unexpected resolved image ref: %+v", doc.ImageRefs[0])
	}
	if doc.InputTitle != "Guide" {
		t.Fatalf("unexpected title: %q", doc.InputTitle)
	}
}
