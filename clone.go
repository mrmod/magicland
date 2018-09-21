package magicland

import (
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func clone(serviceName string, cloneOptions *git.CloneOptions) error {
	repository, err := git.PlainClone("/tmp/"+serviceName+"/stage", false, cloneOptions)
	_ = repository
	if err != nil {
		return err
	}
	return nil
}

// PublicClone Clones a public repository over HTTPS
func PublicClone(gitConfig GitConfiguration) error {
	cloneOptions := &git.CloneOptions{
		URL:           gitConfig.RepositoryURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + gitConfig.BranchName),
		SingleBranch:  true,
		NoCheckout:    false,
	}
	return clone(gitConfig.ServiceName, cloneOptions)
}
