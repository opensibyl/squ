package main

import (
	"flag"

	UnitSqueezer "github.com/opensibyl/UnitSqueezor"
	"github.com/opensibyl/UnitSqueezor/object"
)

func main() {
	// cmd parse
	src := flag.String("src", ".", "repo path")
	before := flag.String("before", "HEAD~1", "before rev")
	after := flag.String("after", "HEAD", "after rev")
	diffOutput := flag.String("diffOutput", "", "diff output")
	dry := flag.Bool("dry", false, "dry")
	flag.Parse()

	config := object.DefaultConfig()
	config.SrcDir = *src
	config.Before = *before
	config.After = *after
	config.DiffFuncOutput = *diffOutput
	config.Dry = *dry

	UnitSqueezer.MainFlow(config)
}
