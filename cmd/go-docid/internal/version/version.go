package version

import "bytes"

var (
	version   = "development"
	goVersion string
	buildTime string
	gitCommit string
	gitTag    string
	gitStatus string
)

// VerbosInfo as name
func VerbosInfo() string {
	var o bytes.Buffer
	o.WriteString("Version:   " + version + "\n")
	o.WriteString("GoVersion: " + goVersion + "\n")
	o.WriteString("BuildTime: " + buildTime + "\n")
	o.WriteString("GitCommit: " + gitCommit + "\n")
	o.WriteString("GitTag:    " + gitTag + "\n")
	o.WriteString("GitStatus: " + gitStatus)
	return o.String()
}

// Version as name
func Version() string {
	return version
}
