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
	diffOutput := flag.String("jsonOutput", "", "diff output")
	dry := flag.Bool("dry", false, "dry")
	indexerType := flag.String("indexer", object.IndexerGolang, "indexer type")
	runnerType := flag.String("runner", object.RunnerGolang, "runner type")
	flag.Parse()

	config := object.DefaultConfig()
	config.SrcDir = *src
	config.Before = *before
	config.After = *after
	config.JsonOutput = *diffOutput
	config.Dry = *dry
	config.IndexerType = *indexerType
	config.RunnerType = *runnerType

	UnitSqueezer.MainFlow(config)
}
