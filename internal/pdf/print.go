package pdf

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/maxanderson95/md2pdf/internal/browser"
	"github.com/maxanderson95/md2pdf/internal/document"
)

type Generator struct {
	Browser browser.Discovery
	Timeout time.Duration
}

type Result struct {
	Warnings []string
}

func (g Generator) Write(ctx context.Context, doc document.Document, outputPath string) (Result, error) {
	timeout := g.Timeout
	if timeout <= 0 {
		timeout = 45 * time.Second
	}

	tmpDir, err := os.MkdirTemp("", "md2pdf-browser-*")
	if err != nil {
		return Result{}, fmt.Errorf("create browser temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpHTML, err := writeTempHTML(doc.HTML)
	if err != nil {
		return Result{}, err
	}
	defer os.Remove(tmpHTML)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(g.Browser.Path),
		chromedp.UserDataDir(filepath.Join(tmpDir, "profile")),
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoFirstRun,
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("metrics-recording-only", true),
	)...)
	defer cancelAlloc()

	runCtx, cancelRun := chromedp.NewContext(allocCtx)
	defer cancelRun()

	runCtx, cancelTimeout := context.WithTimeout(runCtx, timeout)
	defer cancelTimeout()

	pageURL := "file://" + filepath.ToSlash(tmpHTML)
	var failedJSON string
	var pdfBytes []byte

	err = chromedp.Run(runCtx,
		chromedp.Navigate(pageURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(waitForAssetsJS, nil),
		chromedp.Evaluate(failedRemoteAssetsJS, &failedJSON),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBytes = buf
			return nil
		}),
	)
	if err != nil {
		return Result{}, fmt.Errorf("render PDF with browser: %w", err)
	}

	warnings, err := decodeWarnings(failedJSON)
	if err != nil {
		return Result{}, err
	}

	if err := writeAtomically(outputPath, pdfBytes); err != nil {
		return Result{}, err
	}

	return Result{Warnings: warnings}, nil
}

func writeTempHTML(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "md2pdf-*.html")
	if err != nil {
		return "", fmt.Errorf("create temporary HTML file: %w", err)
	}
	defer tmpFile.Close()
	if _, err := tmpFile.WriteString(content); err != nil {
		return "", fmt.Errorf("write temporary HTML file: %w", err)
	}
	return tmpFile.Name(), nil
}

func writeAtomically(outputPath string, contents []byte) error {
	dir := filepath.Dir(outputPath)
	tmpFile, err := os.CreateTemp(dir, "md2pdf-*.pdf")
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	tmpName := tmpFile.Name()
	if _, err := tmpFile.Write(contents); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return fmt.Errorf("write output file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("close output file: %w", err)
	}
	if err := os.Rename(tmpName, outputPath); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("move PDF into place: %w", err)
	}
	return nil
}

func decodeWarnings(raw string) ([]string, error) {
	if raw == "" || raw == "null" {
		return nil, nil
	}
	var warnings []string
	if err := json.Unmarshal([]byte(raw), &warnings); err != nil {
		return nil, fmt.Errorf("decode browser asset warnings: %w", err)
	}
	return warnings, nil
}

const waitForAssetsJS = `new Promise((resolve) => {
  const done = () => resolve(true);
  const ready = () => {
    if (!document.fonts || !document.fonts.ready) {
      return Promise.resolve();
    }
    return document.fonts.ready.catch(() => undefined);
  };
  const images = Array.from(document.images || []);
  Promise.all(images.map((img) => {
    if (img.complete) {
      return Promise.resolve();
    }
    return new Promise((imageResolve) => {
      const finish = () => imageResolve();
      img.addEventListener('load', finish, { once: true });
      img.addEventListener('error', finish, { once: true });
      setTimeout(finish, 10000);
    });
  })).then(() => ready()).then(done).catch(done);
})`

const failedRemoteAssetsJS = `JSON.stringify(Array.from(document.images || [])
  .filter((img) => {
    const src = img.currentSrc || img.src || '';
    return /^https?:\/\//.test(src) && img.complete && img.naturalWidth === 0;
  })
  .map((img) => img.currentSrc || img.src))`
