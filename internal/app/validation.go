package app

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/maxanderson95/md2pdf/internal/document"
)

func validateLocalImages(inputDir string, refs []document.ImageRef) error {
	for _, ref := range refs {
		if ref.IsRemote || ref.IsDataURI {
			continue
		}
		parsed, err := url.Parse(ref.Resolved)
		if err != nil {
			return fmt.Errorf("failed to parse image path %q: %w", ref.Source, err)
		}
		var imagePath string
		switch parsed.Scheme {
		case "", "file":
			imagePath = parsed.Path
		default:
			continue
		}
		if imagePath == "" {
			imagePath = filepath.Join(inputDir, ref.Source)
		}
		if _, err := os.Stat(filepath.Clean(imagePath)); err != nil {
			return fmt.Errorf("local image does not exist: %s", ref.Source)
		}
	}
	return nil
}
