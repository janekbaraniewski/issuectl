package issuectl

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getDefaultConfigFilePath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dirname, ".issuerc")
}

var DefaultConfigFilePath = getDefaultConfigFilePath()

// ProfileName is a name of issuectl config profile
type ProfileName string

// Profile is a config profile
type Profile struct {
	Name       ProfileName       `json:"name"`
	WorkDir    string            `json:"workDir"`
	Repository RepoConfigName    `json:"repository"`
	Backend    BackendConfigName `json:"backend"`
}

// IssuectlConfig manages configuration
type IssuectlConfig struct {
	CurrentProfile    ProfileName     `json:"currentProfile"`
	WorkDir           string          `json:"workDir"`
	DefaultRepository RepoConfigName  `json:"defaultRepository"`
	Repositories      []RepoConfig    `json:"repositories"`
	Issues            []IssueConfig   `json:"issues"`
	Profiles          []Profile       `json:"profiles"`
	Backends          []BackendConfig `json:"backends"`
}

func (c *IssuectlConfig) Save() error {
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		DefaultConfigFilePath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(string(y))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return nil

}

func LoadConfig() *IssuectlConfig {
	config := &IssuectlConfig{}

	data, err := os.ReadFile(DefaultConfigFilePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := config.Save(); err != nil {
				Log.Infof("%v", err)
				return config
			}
			return config
		}
		return nil
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		Log.Infof("%v", err)
		return nil
	}

	return config
}

// Issues

func (ic *IssuectlConfig) AddIssue(issueConfig *IssueConfig) error {
	ic.Issues = append(ic.Issues, *issueConfig)
	return ic.Save()
}

func (ic *IssuectlConfig) DeleteIssue(issueID IssueID) error {
	for i, issueConfig := range ic.Issues {
		if issueConfig.ID == issueID {
			if i < len(ic.Issues) {
				ic.Issues = append(ic.Issues[:i], ic.Issues[i+1:]...)
			} else {
				ic.Issues = ic.Issues[:i]
			}
			return ic.Save()
		}
	}
	return fmt.Errorf("issue with ID '%s' not found", issueID)
}

func (ic *IssuectlConfig) GetIssue(issueID IssueID) *IssueConfig {
	for _, issueConfig := range ic.Issues {
		if issueConfig.ID == issueID {
			return &issueConfig
		}
	}

	return nil
}

// Repositories

func (ic *IssuectlConfig) ListRepositories() error {
	Log.Infof("%v", ic.Repositories)
	return nil
}

func (ic *IssuectlConfig) GetRepository(name RepoConfigName) *RepoConfig {
	for _, rc := range ic.Repositories {
		if rc.Name == name {
			return &rc
		}
	}
	return nil
}

func (ic *IssuectlConfig) AddRepository(repoConfig *RepoConfig) error {
	ic.Repositories = append(ic.Repositories, *repoConfig)
	if err := ic.Save(); err != nil {
		return err
	}
	return nil
}

// Profiles

func (ic *IssuectlConfig) GetProfile(profileName ProfileName) *Profile {
	for _, pr := range ic.Profiles {
		if pr.Name == profileName {
			return &pr
		}
	}
	return nil
}
