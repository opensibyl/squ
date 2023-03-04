package UnitSqueezer

import (
	"testing"

	"github.com/opensibyl/UnitSqueezor/object"
)

func TestGolang(t *testing.T) {
	conf := object.DefaultConfig()
	conf.Dry = true
	MainFlow(conf)
}

func TestMaven(t *testing.T) {
	t.Skip()
	conf := object.DefaultConfig()
	conf.Dry = true
	conf.SrcDir = "../jacoco"
	conf.Before = "HEAD~5"
	conf.IndexerType = object.IndexerJavaJUnit
	conf.RunnerType = object.RunnerMaven
	conf.DiffFuncOutput = "./output.json"
	MainFlow(conf)
}
