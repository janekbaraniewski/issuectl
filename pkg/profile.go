package issuectl

// ProfileName is a name of issuectl config profile
type ProfileName string

// Profile is a config profile
type Profile struct {
	Name         ProfileName       `yaml:"name"`
	WorkDir      string            `yaml:"workDir"`
	IssueBackend BackendConfigName `yaml:"issueBackend"`
	RepoBackend  BackendConfigName `yaml:"repoBackend"`
	GitUserName  GitUserName       `yaml:"gituser"`
	Repositories []RepoConfigName  `yaml:"repositories"`

	// DefaultRepository is now used for Github IssueBackend
	DefaultRepository RepoConfigName `yaml:"defaultRepository"`
}

func (p *Profile) AddRepository(repo RepoConfigName) error {
	p.Repositories = append(p.Repositories, repo)
	return nil
}
