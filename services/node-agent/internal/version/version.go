package version

type Info struct {
	Name    string
	Version string
	Commit  string
	BuiltAt string
}

var Build = Info{
	Name:    "Sentinel Node Agent",
	Version: "0.3.0-dev",
	Commit:  "development",
	BuiltAt: "unknown",
}
