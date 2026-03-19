package version

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func Info() BuildInfo {
	return BuildInfo{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

func (b BuildInfo) String() string {
	return fmt.Sprintf("md2pdf %s (commit %s, built %s)", b.Version, b.Commit, b.Date)
}
