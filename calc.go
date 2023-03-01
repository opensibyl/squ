package UnitSqueezer

import (
	"os/exec"

	"github.com/opensibyl/UnitSqueezor/object"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/ext"
)

type DiffMap = map[string][]int

type DiffExtractor interface {
	ExtractDiffMap() (DiffMap, error)
}

func NewDiffExtractor(config *object.SharedConfig) (DiffExtractor, error) {
	return &GitExtractor{config}, nil
}

type GitExtractor struct {
	config *object.SharedConfig
}

func (g *GitExtractor) ExtractDiffMap() (DiffMap, error) {
	gitDiffCmd := exec.Command("git", "diff", "HEAD~1", "HEAD")
	gitDiffCmd.Dir = g.config.SrcDir
	patchRaw, err := gitDiffCmd.CombinedOutput()
	if err != nil {
		core.Log.Errorf("git cmd error: %s", patchRaw)
		panic(err)
	}

	return ext.Unified2Affected(patchRaw)
}
