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

func AddIssue(issueConfig *IssueConfig) error {
	config := LoadConfig()

	config.Issues = append(config.Issues, *issueConfig)

	return config.Save()
}

func DeleteIssue(issueID IssueID) error {
	config := LoadConfig()

	for i, ic := range config.Issues {
		if ic.ID == issueID {
			if i < len(config.Issues) {
				config.Issues = append(config.Issues[:i], config.Issues[i+1:]...)
			} else {
				config.Issues = config.Issues[:i]
			}
			config.Save()
			return nil
		}
	}

	return nil
}

func GetIssue(issueID IssueID) *IssueConfig {
	config := LoadConfig()

	for _, ic := range config.Issues {
		if ic.ID == issueID {
			return &ic
		}
	}

	return nil
}

// Repositories

func ListRepositories() error {
	config := LoadConfig()

	Log.Infof("%v", config.Repositories)

	return nil
}

func GetRepository(name RepoConfigName) *RepoConfig {
	config := LoadConfig()

	for _, rc := range config.Repositories {
		if rc.Name == name {
			return &rc
		}
	}

	return nil
}

func AddRepository(repoConfig *RepoConfig) error {
	config := LoadConfig()

	config.Repositories = append(config.Repositories, *repoConfig)

	if err := config.Save(); err != nil {
		Log.Infof("ERROR!!!! %v", err)
		return err
	}

	return nil
}
