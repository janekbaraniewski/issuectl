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
