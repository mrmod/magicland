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

func (this GitConfiguration) servicePath() string {
	return "/tmp/" + this.ServiceName + "/stage/"
}

// AsRedisHashMap Convert to a storable datastructure
func (this GitConfiguration) AsRedisHashMap() map[string]interface{} {
	return map[string]interface{}{
		"user":          this.User,
		"branchname":    this.BranchName,
		"repositoryurl": this.RepositoryURL,
		"servicename":   this.ServiceName,
	}
}

// FromRedisHashMap Create a new configuration
func (this GitConfiguration) FromRedisHashMap(hm []interface{}) GitConfiguration {
	return GitConfiguration{
		User:          hm[0].(string),
		BranchName:    hm[1].(string),
		RepositoryURL: hm[2].(string),
		ServiceName:   hm[3].(string),
	}
}
