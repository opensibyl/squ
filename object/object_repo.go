package object

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

type RepoInfo struct {
	RepoId  string `json:"repoId"`
	RevHash string `json:"revHash"`
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
		RepoId:  filepath.Base(srcDir),
		RevHash: head.Hash().String(),
	}, nil
}
