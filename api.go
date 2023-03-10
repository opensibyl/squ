package UnitSqueezer

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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

	rootContext := context.Background()
	sibylContext, cancel := context.WithCancel(rootContext)
	defer cancel()

	rootStartTime := time.Now()

	// 0. start sibyl2 backend
	sibylContext, stop := signal.NotifyContext(rootContext, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		config := object2.DefaultExecuteConfig()
		// for performance
		config.BindingConfigPart.DbType = object2.DriverTypeInMemory
		config.EnableLog = false
		err := server.Execute(config, sibylContext)
		PanicIfErr(err)
	}()
	defer stop()
	log.Log.Infof("sibyl2 backend ready")

	// 1. upload and tag
	curIndexer, err := indexer.GetIndexer(conf.IndexerType, &conf)
	err = curIndexer.UploadSrc(sibylContext)
	PanicIfErr(err)
	err = curIndexer.TagCases(sibylContext)
	PanicIfErr(err)
	log.Log.Infof("indexer ready")

	// 2. calc
	// line level diff
	calcContext, cancel := context.WithCancel(rootContext)
	defer cancel()
	curExtractor, err := extractor.NewDiffExtractor(&conf)
	PanicIfErr(err)
	diffMap, err := curExtractor.ExtractDiffMethods(calcContext)
	PanicIfErr(err)
	log.Log.Infof("diff calc ready: %v", len(diffMap))

	// 3. runner
	log.Log.Infof("runner scope")
	runnerContext, cancel := context.WithCancel(rootContext)
	defer cancel()
	curRunner, err := runner.GetRunner(conf.RunnerType, conf)
	PanicIfErr(err)

	cache := curIndexer.GetTagCache()
	casesToRun := make([]*openapi.ObjectFunctionWithSignature, 0)
	for fileName, eachFunctionList := range diffMap {
		log.Log.Infof("handle modified file: %s, functions: %d", fileName, len(eachFunctionList))
		for _, eachFunc := range eachFunctionList {
			for eachCaseSignature, eachCase := range cache {
				log.Log.Infof("case %v influence %v", eachCaseSignature, len(eachCaseSignature))
				if _, ok := (*eachCase)[eachFunc.GetSignature()]; ok {
					// related case
					log.Log.Infof("case %v related to %v", eachCaseSignature, eachFunc.GetSignature())
				} else {
					// not related
				}
			}
		}
	}
	PanicIfErr(err)
	log.Log.Infof("case analyzer done")

	if conf.JsonOutput != "" {
		o := &Output{}
		o.DiffMap = diffMap
		o.Cases = casesToRun
		diffFuncOutputBytes, err := json.Marshal(o)
		PanicIfErr(err)
		err = os.WriteFile(conf.JsonOutput, diffFuncOutputBytes, os.ModePerm)
		PanicIfErr(err)
	}

	// shutdown sibyl2 before running cases because of its ports
	stop()
	prepareTotalCost := time.Since(rootStartTime)
	log.Log.Infof("prepare stage finished, total cost: %d ms", prepareTotalCost.Milliseconds())
	if !conf.Dry {
		log.Log.Infof("start running cases: %v", len(casesToRun))
		err = curRunner.Run(casesToRun, runnerContext)
		PanicIfErr(err)
	}
}

type Output struct {
	DiffMap object.DiffFuncMap                     `json:"diff"`
	Cases   []*openapi.ObjectFunctionWithSignature `json:"cases"`
}
