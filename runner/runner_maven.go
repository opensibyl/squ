package runner

import (
	"context"
	"fmt"

	"github.com/opensibyl/UnitSqueezor/log"
	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type MavenRunner struct {
	*BaseRunner
}

func (m *MavenRunner) Run(cases []*openapi.ObjectFunctionWithSignature, _ context.Context) error {
	// mvn test -Dtest="TheSecondUnitTest#whenTestCase2_thenPrintTest2_1"
	parts := make([]string, 0, len(cases))
	for _, each := range cases {
		extras := each.GetExtras()
		log.Log.Infof("map: %v", extras)

		clazzName, clazzNameExisted := extras["className"].(string)
		if !clazzNameExisted {
			log.Log.Warnf("class name not found in: %v", each.GetName())
			continue
		}
		curPartStr := fmt.Sprintf("%s#%s", clazzName, each.GetName())
		parts = append(parts, curPartStr)
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
