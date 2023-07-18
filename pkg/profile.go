package issuectl

// ProfileName is a name of issuectl config profile
type ProfileName string

// Profile is a config profile
type Profile struct {
	Name         ProfileName       `json:"name"`
	WorkDir      string            `json:"workDir"`
	Repository   RepoConfigName    `json:"repository"`
	Backend      BackendConfigName `json:"backend"`
	Repositories []*RepoConfig     `json:"repositories"`
}

func (p *Profile) AddRepository(repo *RepoConfig) error {
	p.Repositories = append(p.Repositories, repo)
	return nil
}
