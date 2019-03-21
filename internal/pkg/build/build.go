package build

import "fmt"

type BuildInfo struct {
	Version   string
	BuildDate string
}

func (b *BuildInfo) String() string {
	v := b.Version
	d := b.BuildDate
	if b.Version == "" {
		v = "UNKNOWN"
	}
	if b.BuildDate == "" {
		d = "UNKNOWN"
	}
	return fmt.Sprintf("%s built on %s", v, d)
}
