package squ

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/pkg/server"
	object2 "github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/squ/extractor"
	"github.com/opensibyl/squ/indexer"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/mapper"
	"github.com/opensibyl/squ/object"
	"github.com/opensibyl/squ/runner"
)

/*
index (preparation)
1. upload all the files to sibyl2
2. search and tag all the test methods
3. calc and tag all the test methods influencing scope

extract (extract data from workspace)
1. calc diff between current and previous
2. find methods influenced by diff

mapper (mapping between cases and diff)
1. search related cases

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

	rootContext := context.Background()
	rootContext, cancel := context.WithCancel(rootContext)
	defer cancel()

	rootStartTime := time.Now()

	// 0. start sibyl2 backend
	if conf.LocalSibyl() {
		log.Log.Infof("using local sibyl, starting ...")
		go func() {
			config := object2.DefaultExecuteConfig()
			// for performance
			config.BindingConfigPart.DbType = object2.DriverTypeInMemory
			config.EnableLog = false
			err := server.Execute(config, rootContext)
			PanicIfErr(err)
		}()
		log.Log.Infof("sibyl2 backend ready")
	} else {
		log.Log.Infof("using remote sibyl, skip starting server")
	}

	// 1. index
	// todo: if using remote sibyl server, pull data directly
	curIndexer, err := indexer.GetIndexer(conf.IndexerType, &conf)
	err = curIndexer.UploadSrc(rootContext)
	PanicIfErr(err)
	err = curIndexer.TagCases(rootContext)
	PanicIfErr(err)
	log.Log.Infof("indexer ready")

	// 2. extract
	// line level diff
	calcContext, cancel := context.WithCancel(rootContext)
	defer cancel()
	curExtractor, err := extractor.NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := curExtractor.ExtractDiffMethods(calcContext)
	PanicIfErr(err)
	log.Log.Infof("diff calc ready: %v", len(diffMap))

	// 3. mapper
	curMapper := mapper.NewMapper()
	curMapper.SetIndexer(curIndexer)
	casesToRun, err := curMapper.Diff2Cases(rootContext, diffMap)
	PanicIfErr(err)
	log.Log.Infof("case analyzer done, before: %d, after: %d", len(curIndexer.GetCaseSet()), len(casesToRun))

	if conf.JsonOutput != "" {
		o := &Output{}
		o.DiffMap = diffMap
		o.Cases = casesToRun
		diffFuncOutputBytes, err := json.Marshal(o)
		PanicIfErr(err)
		err = os.WriteFile(conf.JsonOutput, diffFuncOutputBytes, os.ModePerm)
		PanicIfErr(err)
	}

	// 4. runner
	log.Log.Infof("runner scope")
	curRunner, err := runner.GetRunner(conf.RunnerType, conf)
	PanicIfErr(err)

	prepareTotalCost := time.Since(rootStartTime)
	log.Log.Infof("prepare stage finished, total cost: %d ms", prepareTotalCost.Milliseconds())
	if conf.Dry {
		cmd := curRunner.GetRunCommand(casesToRun)
		log.Log.Infof("runner cmd: %v", cmd)
	} else {
		log.Log.Infof("start running cases: %v", len(casesToRun))
		err = curRunner.Run(casesToRun, rootContext)
		PanicIfErr(err)
	}
}

type Output struct {
	DiffMap object.DiffFuncMap                     `json:"diff"`
	Cases   []*openapi.ObjectFunctionWithSignature `json:"cases"`
}
