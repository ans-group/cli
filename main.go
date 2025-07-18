package main

import (
	"github.com/ans-group/cli/cmd"
	"github.com/ans-group/cli/internal/pkg/build"
)

func main() {
	cmd.Execute(build.BuildInfo{Version: VERSION, BuildDate: BUILDDATE})
}
