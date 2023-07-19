package issuectl

// ProfileName is a name of issuectl config profile
type ProfileName string

// Profile is a config profile
type Profile struct {
	Name         ProfileName       `json:"name"`
	WorkDir      string            `json:"workDir"`
	Backend      BackendConfigName `json:"backend"`
	GitUserName  GitUserName       `json:"gituser"`
	Repositories []*RepoConfig     `json:"repositories"`

	// DefaultRepository is now used for Github IssueBackend
	DefaultRepository RepoConfigName `json:"defaultRepository"`
}

func (p *Profile) AddRepository(repo *RepoConfig) error {
	p.Repositories = append(p.Repositories, repo)
	return nil
}
