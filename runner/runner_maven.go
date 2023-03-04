package runner

import (
	"context"
	"strings"

	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type MavenRunner struct {
	*BaseRunner
}

func (m *MavenRunner) Run(cases []*openapi.ObjectFunctionWithSignature, ctx context.Context) error {
	// mvn test -Dtest="TheSecondUnitTest#whenTestCase2_thenPrintTest2_1"
	parts := make([]string, 0, len(cases))
	for _, each := range cases {
		extras := each.GetExtras()
		curPart := strings.Join([]string{
			extras["packageName"].(string),
			extras["className"].(string),
			each.GetName()}, "")
		parts = append(parts, curPart)
	}
	log.Log.Infof("%v", parts)
	return nil
}

func NewMavenRunner(conf *object.SharedConfig) (Runner, error) {
	apiClient, err := conf.NewSibylClient()
	if err != nil {
		return nil, err
	}
	return &MavenRunner{
		&BaseRunner{
			conf,
			apiClient,
		},
	}, nil
}
