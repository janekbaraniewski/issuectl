package issuectl

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var configFilePath = "/Users/janbaraniewski/.issuerc" // FIXME

// ProfileName is a name of issuectl config profile
type ProfileName string

// Profile is a config profile
type Profile struct {
	Name              ProfileName    `json:"name"`
	WorkDir           string         `json:"workDir"`
	DefaultRepository RepoConfigName `json:"defaultRepository"`
}

// IssuectlConfig manages configuration
type IssuectlConfig struct {
	CurrentProfile    ProfileName    `json:"currentProfile"`
	WorkDir           string         `json:"workDir"`
	DefaultRepository RepoConfigName `json:"defaultRepository"`
	Repositories      []RepoConfig   `json:"repositories"`
	Issues            []IssueConfig  `json:"issues"`
	Profiles          []Profile      `json:"profiles"`
}

func (c *IssuectlConfig) Save() error {
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		configFilePath,
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

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		Log.Infof("%v", err)
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
			ic.Save()
			return nil
		}
	}
	return nil
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
