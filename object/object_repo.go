package object

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

type RepoInfo struct {
	Name     string `json:"name"`
	CommitId string `json:"commitId"`
}

func GetRepoInfoFromDir(srcDir string) (*RepoInfo, error) {
	repo, err := git.PlainOpen(srcDir)
	if err != nil {
		return nil, fmt.Errorf("load repo from %s failed", srcDir)
	}
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &RepoInfo{
		Name:     filepath.Base(srcDir),
		CommitId: head.Hash().String(),
	}, nil
}
