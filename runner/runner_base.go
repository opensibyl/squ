package runner

import (
	"context"

	"github.com/opensibyl/UnitSqueezor/object"
	openapi "github.com/opensibyl/sibyl-go-client"
)

type BaseRunner struct {
	config    *object.SharedConfig
	apiClient *openapi.APIClient
}

func (baseRunner *BaseRunner) GetRelatedCases(ctx context.Context, caseTagCache object.CaseTagCache, targetSignature string) ([]*openapi.ObjectFunctionWithSignature, error) {
	relatedCases := make([]string, 0)
	for caseSignature, m := range caseTagCache {
		if _, ok := m[targetSignature]; ok {
			relatedCases = append(relatedCases, caseSignature)
		}
	}

	r := make([]*openapi.ObjectFunctionWithSignature, 0, len(relatedCases))
	for _, each := range relatedCases {

		cur, _, err := baseRunner.apiClient.
			SignatureQueryApi.
			ApiV1SignatureFuncGet(ctx).
			Repo(baseRunner.config.RepoInfo.RepoId).
			Rev(baseRunner.config.RepoInfo.RevHash).
			Signature(each).
			Execute()
		if err != nil {
			return nil, err
		}
		r = append(r, cur)
	}
	return r, nil
}
