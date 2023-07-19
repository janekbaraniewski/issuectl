package issuectl

import (
	"errors"
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

func getDefaultSSHKeyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".ssh/id_ed25519")
}

var DefaultConfigFilePath = getDefaultConfigFilePath()
var DefaultSSHKeyPath = getDefaultSSHKeyPath()

// IssuectlConfig manages configuration
type IssuectlConfig struct {
	CurrentProfile ProfileName                         `json:"currentProfile"`
	Repositories   map[RepoConfigName]RepoConfig       `json:"repositories"`
	Issues         map[IssueID]IssueConfig             `json:"issues"`
	Profiles       map[ProfileName]Profile             `json:"profiles"`
	Backends       map[BackendConfigName]BackendConfig `json:"backends"`
	GitUsers       map[GitUserName]GitUser             `json:"gitUsers"`
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
	ic.Issues[issueConfig.ID] = *issueConfig
	return ic.Save()
}

func (ic *IssuectlConfig) DeleteIssue(issueID IssueID) error {
	delete(ic.Issues, issueID)
	return ic.Save()
}

func (ic *IssuectlConfig) GetIssue(issueID IssueID) (IssueConfig, bool) {
	issue, ok := ic.Issues[issueID]
	return issue, ok
}

// Repositories

func (ic *IssuectlConfig) ListRepositories() error {
	Log.Infof("%v", ic.Repositories)
	return nil
}

func (ic *IssuectlConfig) GetRepository(name RepoConfigName) RepoConfig {
	return ic.Repositories[name]
}

func (ic *IssuectlConfig) AddRepository(repoConfig *RepoConfig) error {
	ic.Repositories[repoConfig.Name] = *repoConfig
	return ic.Save()
}

// Profiles

func (ic *IssuectlConfig) GetProfile(profileName ProfileName) Profile {
	return ic.Profiles[profileName]
}

func (ic *IssuectlConfig) AddProfile(profile *Profile) error {
	ic.Profiles[profile.Name] = *profile
	return ic.Save()
}

func (ic *IssuectlConfig) DeleteProfile(profileName ProfileName) error {
	delete(ic.Profiles, profileName)
	return ic.Save()
}

func (ic *IssuectlConfig) GetCurrentProfile() ProfileName {
	return ic.CurrentProfile
}

func (ic *IssuectlConfig) UseProfile(profile ProfileName) error {
	ic.CurrentProfile = profile
	return ic.Save()
}

func (ic *IssuectlConfig) GetProfiles() map[ProfileName]Profile {
	return ic.Profiles
}

func (ic *IssuectlConfig) UpdateProfile(profile *Profile) error {
	ic.Profiles[profile.Name] = *profile
	return ic.Save()
}

// Backends
func (ic *IssuectlConfig) GetBackend(backendName BackendConfigName) BackendConfig {
	return ic.Backends[backendName]
}

func (ic *IssuectlConfig) AddBackend(backend *BackendConfig) error {
	ic.Backends[backend.Name] = *backend
	return ic.Save()
}

func (ic *IssuectlConfig) DeleteBackend(backendName BackendConfigName) error {
	delete(ic.Backends, backendName)
	return ic.Save()
}

func (ic *IssuectlConfig) GetBackends() map[BackendConfigName]BackendConfig {
	return ic.Backends
}

func (ic *IssuectlConfig) GetRepositories() map[RepoConfigName]RepoConfig {
	return ic.Repositories
}

// GitUsers
func (ic *IssuectlConfig) GetGitUser(userName GitUserName) (GitUser, bool) {
	gitUser, exists := ic.GitUsers[userName]
	return gitUser, exists
}

func (ic *IssuectlConfig) AddGitUser(user *GitUser) error {
	ic.GitUsers[GitUserName(user.Name)] = *user
	return ic.Save()
}

func (ic *IssuectlConfig) DeleteGitUser(userName GitUserName) error {
	delete(ic.GitUsers, userName)
	return ic.Save()
}

func (ic *IssuectlConfig) GetGitUsers() map[GitUserName]GitUser {
	return ic.GitUsers
}
