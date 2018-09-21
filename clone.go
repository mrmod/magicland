package magicland

import (
	git "gopkg.in/src-d/go-git.v4"
)

func clone(serviceName string, cloneOptions *git.CloneOptions) error {
	repository, err := git.PlainClone("/tmp/"+serviceName+"/stage", false, cloneOptions)
	_ = repository
	if err != nil {
		return err
	}
	return nil
}

func PublicClone(gitConfig GitConfiguration) error {
	cloneOptions := &git.CloneOptions{
		URL: gitConfig.RepositoryURL,
		// ReferenceName: plumbing.ReferenceName(gitConfig.BranchName),
		SingleBranch: true,
		NoCheckout:   false,
	}
	return clone(gitConfig.ServiceName, cloneOptions)
}
