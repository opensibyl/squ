package UnitSqueezer

import (
	"context"
	"path/filepath"

	"github.com/opensibyl/UnitSqueezor/indexer"
	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

/*
index
1. upload all the files to sibyl2
2. search and tag all the test methods
3. calc and tag all the test methods influencing scope

calc
1. calc diff between current and previous
2. find methods influenced by diff
3. search related cases

execute (can be implemented by different languages
1. build test commands for different frameworks
2. call cmd
*/

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MainFlow() {
	sharedContext := context.Background()

	// init config
	conf := object.DefaultConfig()
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir

	// todo: start sibyl2 server

	// 1. upload and tag
	curIndexer, err := indexer.NewIndexer(&conf)
	err = curIndexer.UploadSrc(sharedContext)
	PanicIfErr(err)
	err = curIndexer.TagCases(sharedContext)
	PanicIfErr(err)

	// 2. calc
	// line level diff
	extractor, err := NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := extractor.ExtractDiffMethods(sharedContext)
	PanicIfErr(err)
	log.Log.Infof("diff map: %v", diffMap)

	// 3. executor
	executor, err := NewGoExecutor(&conf)
	PanicIfErr(err)
	cases, err := executor.GetRelatedCases(sharedContext, diffMap)
	PanicIfErr(err)
	err = executor.Execute(cases, sharedContext)
	PanicIfErr(err)
}

type FunctionWithState struct {
	*openapi.ObjectFunctionWithSignature

	Reachable bool
	ReachBy   []string
}
