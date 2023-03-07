package UnitSqueezer

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/opensibyl/UnitSqueezor/extractor"
	"github.com/opensibyl/UnitSqueezor/indexer"
	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	"github.com/opensibyl/UnitSqueezor/runner"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/pkg/server"
	object2 "github.com/opensibyl/sibyl2/pkg/server/object"
)

/*
index
1. upload all the files to sibyl2
2. search and tag all the test methods
3. calc and tag all the test methods influencing scope

extract
1. calc diff between current and previous
2. find methods influenced by diff
3. search related cases

runner (can be implemented by different languages)
1. build test commands for different frameworks
2. call cmd
*/

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MainFlow(conf object.SharedConfig) {
	// init config
	absSrcDir, err := filepath.Abs(conf.SrcDir)
	PanicIfErr(err)
	conf.SrcDir = absSrcDir

	// 0. start sibyl2 backend
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		config := object2.DefaultExecuteConfig()
		// for performance
		config.BindingConfigPart.DbType = object2.DriverTypeInMemory
		err := server.Execute(config, ctx)
		PanicIfErr(err)
	}()
	defer stop()
	log.Log.Infof("sibyl2 backend ready")

	// 1. upload and tag
	curIndexer, err := indexer.GetIndexer(conf.IndexerType, &conf)
	err = curIndexer.UploadSrc(ctx)
	PanicIfErr(err)
	err = curIndexer.TagCases(ctx)
	PanicIfErr(err)
	log.Log.Infof("indexer ready")

	// 2. calc
	// line level diff
	curExtractor, err := extractor.NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := curExtractor.ExtractDiffMethods(ctx)
	PanicIfErr(err)
	log.Log.Infof("diff calc ready: %v", len(diffMap))
	if conf.DiffFuncOutput != "" {
		diffFuncOutputBytes, err := json.Marshal(diffMap)
		PanicIfErr(err)
		err = os.WriteFile(conf.DiffFuncOutput, diffFuncOutputBytes, os.ModePerm)
		PanicIfErr(err)
	}

	// 3. runner
	tagCache := curIndexer.GetTagMap()
	curRunner, err := runner.GetRunner(conf.RunnerType, conf)
	PanicIfErr(err)

	caseSet := make(map[string]*openapi.ObjectFunctionWithSignature)
	for _, eachFunctionList := range diffMap {
		for _, eachFunc := range eachFunctionList {
			cases, err := curRunner.GetRelatedCases(ctx, *tagCache, eachFunc.GetSignature())
			PanicIfErr(err)
			for _, eachCase := range cases {
				caseSet[eachCase.GetSignature()] = eachCase
			}
		}
	}
	PanicIfErr(err)

	casesToRun := make([]*openapi.ObjectFunctionWithSignature, 0, len(caseSet))
	for _, each := range caseSet {
		casesToRun = append(casesToRun, each)
	}

	if !conf.Dry {
		// shutdown sibyl2 before running cases because of its ports
		stop()
		log.Log.Infof("start running cases: %v", len(caseSet))

		runnerContext := context.Background()
		err = curRunner.Run(casesToRun, runnerContext)
		PanicIfErr(err)
	}
}
