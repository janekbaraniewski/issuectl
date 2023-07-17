package issuectl

var workDir string = "/Users/janbaraniewski/Workspace/priv/issuectl/testWorkdir" // FIXME

// IssuectlConfig manages configuration
type IssuectlConfig struct {
	WorkDir           string         `json:"workDir"`
	DefaultRepository RepoConfigName `json:"defaultRepository"`
	Repositories      []RepoConfig   `json:"repositories"`
}

func (c *IssuectlConfig) Save() error {
	return nil
}

func LoadConfig() *IssuectlConfig {
	return &IssuectlConfig{
		WorkDir:           workDir,
		DefaultRepository: "multi-cloud",
		Repositories: []RepoConfig{
			{
				Name:    "multi-cloud",
				RepoUrl: "git@github.com:elotl/multi-cloud.git",
			},
		},
	}
}
