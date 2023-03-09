package UnitSqueezer

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
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

	rootContext := context.Background()
	sibylContext, cancel := context.WithCancel(rootContext)
	defer cancel()

	// 0. start sibyl2 backend
	sibylContext, stop := signal.NotifyContext(sibylContext, syscall.SIGINT, syscall.SIGTERM)
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
	if conf.DiffFuncOutput != "" {
		diffFuncOutputBytes, err := json.Marshal(diffMap)
		PanicIfErr(err)
		err = os.WriteFile(conf.DiffFuncOutput, diffFuncOutputBytes, os.ModePerm)
		PanicIfErr(err)
	}

	// 3. runner
	log.Log.Infof("runner scope")
	runnerContext, cancel := context.WithCancel(rootContext)
	defer cancel()
	curRunner, err := runner.GetRunner(conf.RunnerType, conf)
	PanicIfErr(err)

	caseSet := make(map[string]*openapi.ObjectFunctionWithSignature)
	for fileName, eachFunctionList := range diffMap {
		log.Log.Infof("handle modified file: %s", fileName)
		var wg sync.WaitGroup
		wg.Add(len(eachFunctionList))
		for _, eachFunc := range eachFunctionList {
			go func(f object.FunctionWithState) {
				defer wg.Done()
				log.Log.Infof("handle modified func: %v", f.GetSignature())
				cases, err := curRunner.GetRelatedCases(runnerContext, f.GetSignature())
				PanicIfErr(err)

				var l sync.Mutex
				l.Lock()
				defer l.Unlock()
				for _, eachCase := range cases {
					caseSet[eachCase.GetSignature()] = eachCase
				}
			}(*eachFunc)
		}
		wg.Wait()
	}
	PanicIfErr(err)

	casesToRun := make([]*openapi.ObjectFunctionWithSignature, 0, len(caseSet))
	for _, each := range caseSet {
		casesToRun = append(casesToRun, each)
	}

	// shutdown sibyl2 before running cases because of its ports
	stop()
	if !conf.Dry {
		log.Log.Infof("start running cases: %v", len(caseSet))
		err = curRunner.Run(casesToRun, runnerContext)
		PanicIfErr(err)
	}
}
