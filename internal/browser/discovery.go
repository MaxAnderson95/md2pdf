package browser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const EnvBrowser = "MD2PDF_BROWSER"

type Discovery struct {
	Path    string
	Tried   []string
	Source  string
	Browser string
}

func Find() (Discovery, error) {
	tried := make([]string, 0)
	if override := strings.TrimSpace(os.Getenv(EnvBrowser)); override != "" {
		tried = append(tried, override)
		if info, err := os.Stat(override); err == nil && !info.IsDir() {
			return Discovery{Path: override, Tried: tried, Source: EnvBrowser, Browser: filepath.Base(override)}, nil
		}
		return Discovery{}, fmt.Errorf("browser override %s does not point to a valid executable", override)
	}

	for _, candidate := range candidatePaths() {
		tried = append(tried, candidate)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return Discovery{Path: candidate, Tried: tried, Source: "well-known path", Browser: filepath.Base(candidate)}, nil
		}
	}

	for _, name := range pathExecutables() {
		tried = append(tried, name)
		if resolved, err := exec.LookPath(name); err == nil {
			return Discovery{Path: resolved, Tried: tried, Source: "PATH", Browser: name}, nil
		}
	}

	return Discovery{}, fmt.Errorf("no supported browser found; tried %s; set %s to override", strings.Join(tried, ", "), EnvBrowser)
}

func candidatePaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		}
	case "windows":
		pf := os.Getenv("ProgramFiles")
		x86 := os.Getenv("ProgramFiles(x86)")
		local := os.Getenv("LocalAppData")
		return []string{
			filepath.Join(pf, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(x86, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(local, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(pf, "Microsoft", "Edge", "Application", "msedge.exe"),
			filepath.Join(x86, "Microsoft", "Edge", "Application", "msedge.exe"),
		}
	default:
		return nil
	}
}

func pathExecutables() []string {
	return []string{"google-chrome", "chromium", "chromium-browser", "microsoft-edge", "chrome", "msedge"}
}
