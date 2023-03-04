package UnitSqueezer

import (
	"bytes"
	"context"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/opensibyl/UnitSqueezor/extractor"
	"github.com/opensibyl/UnitSqueezor/indexer"
	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	"github.com/opensibyl/UnitSqueezor/runner"
	server2 "github.com/opensibyl/sibyl2/cmd/sibyl/subs/server"
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

func MainFlow(conf object.SharedConfig) {
	sharedContext := context.Background()

	// init config
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir

	// 0. start sibyl2 backend
	cmd := server2.NewServerCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go cmd.ExecuteContext(ctx)
	defer stop()
	log.Log.Infof("sibyl2 backend ready")

	// 1. upload and tag
	curIndexer, err := indexer.NewIndexer(&conf)
	err = curIndexer.UploadSrc(sharedContext)
	PanicIfErr(err)
	err = curIndexer.TagCases(sharedContext)
	PanicIfErr(err)
	log.Log.Infof("indexer ready")

	// 2. calc
	// line level diff
	curExtractor, err := extractor.NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := curExtractor.ExtractDiffMethods(sharedContext)
	PanicIfErr(err)
	log.Log.Infof("diff calc ready: %v", len(diffMap))

	// 3. executor
	executor, err := runner.NewGolangRunner(&conf)
	PanicIfErr(err)
	cases, err := executor.GetRelatedCases(sharedContext, diffMap)
	PanicIfErr(err)
	log.Log.Infof("start running cases: %v", len(cases))
	err = executor.Run(cases, sharedContext)
	PanicIfErr(err)
}
