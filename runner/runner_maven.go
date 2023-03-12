package runner

import (
	"fmt"
	"strings"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/object"
)

type MavenRunner struct {
	*BaseRunner
}

func (m *MavenRunner) GetRunCommand(cases []*openapi.ObjectFunctionWithSignature) []string {
	// mvn test -Dtest="TheSecondUnitTest#whenTestCase2_thenPrintTest2_1"
	parts := make([]string, 0, len(cases))
	for _, each := range cases {
		extras := each.GetExtras()
		log.Log.Infof("map: %v", extras)

		clazzInfo, existed := extras["classInfo"].(map[string]interface{})
		if !existed {
			log.Log.Warnf("class info not found in: %v", each.GetName())
			continue
		}
		packageName, packageExisted := clazzInfo["packageName"].(string)
		if !packageExisted {
			log.Log.Warnf("package name not found in: %v", each.GetName())
			continue
		}
		clazzName, nameExisted := clazzInfo["className"].(string)
		if !nameExisted {
			log.Log.Warnf("class name not found in: %v", each.GetName())
			continue
		}

		curPartStr := fmt.Sprintf("%s#%s", fmt.Sprintf("%s.%s", packageName, clazzName), each.GetName())
		parts = append(parts, curPartStr)
	}
	joined := strings.Join(parts, ",")
	return []string{"mvn", "test", "-Dtest=" + joined, "-DfailIfNoTests=false"}
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
