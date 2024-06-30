package main

import (
	"github.com/mrkovshik/yametrics/cmd/staticlint/analizer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analizer.OSExitAnalyzer)
}
