package issuectl

// IssuectlConfig manages configuration
type IssuectlConfig struct {
	Repositories []RepoConfig `json:"repositories"`
}

func (c *IssuectlConfig) Save() error {
	return nil
}
