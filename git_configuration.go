package magicland

// GitConfiguration For a user of the service, there are
// a few things to keep track of in order to clone/checkout
// their git repository
type GitConfiguration struct {
	User          string
	BranchName    string
	RepositoryURL string
	ServiceName   string
}

func NewGitConfiguration(user, branch, url, service string) GitConfiguration {
	return GitConfiguration{
		User:          user,
		BranchName:    branch,
		RepositoryURL: url,
		ServiceName:   service,
	}
}
