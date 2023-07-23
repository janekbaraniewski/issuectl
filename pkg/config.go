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
type issuectlConfig struct {
	CurrentProfile ProfileName                          `yaml:"currentProfile"`
	Repositories   map[RepoConfigName]*RepoConfig       `yaml:"repositories,omitempty"`
	Issues         map[IssueID]*IssueConfig             `yaml:"issues,omitempty"`
	Profiles       map[ProfileName]*Profile             `yaml:"profiles,omitempty"`
	Backends       map[BackendConfigName]*BackendConfig `yaml:"backends,omitempty"`
	GitUsers       map[GitUserName]*GitUser             `yaml:"gitUsers,omitempty"`

	_persistenceMode string `yaml:"-"`
}

type IssuectlConfig interface {
	GetInMemory() IssuectlConfig
	GetPersistent() IssuectlConfig

	// Profile
	AddProfile(*Profile) error
	DeleteProfile(profileName ProfileName) error
	GetCurrentProfile() ProfileName
	GetProfile(ProfileName) *Profile
	GetProfiles() map[ProfileName]*Profile
	UpdateProfile(*Profile) error
	UseProfile(profile ProfileName) error

	// Issues
	AddIssue(issueConfig *IssueConfig) error
	DeleteIssue(issueID IssueID) error
	GetIssue(IssueID) (*IssueConfig, bool)
	GetIssues() map[IssueID]*IssueConfig

	// Repositories
	AddRepository(repoConfig *RepoConfig) error
	GetRepository(name RepoConfigName) *RepoConfig
	GetRepositories() map[RepoConfigName]*RepoConfig

	// Backends
	AddBackend(backend *BackendConfig) error
	DeleteBackend(backendName BackendConfigName) error
	GetBackend(backendName BackendConfigName) *BackendConfig
	GetBackends() map[BackendConfigName]*BackendConfig

	// GitUsers
	AddGitUser(user *GitUser) error
	DeleteGitUser(userName GitUserName) error
	GetGitUser(userName GitUserName) (*GitUser, bool)
	GetGitUsers() map[GitUserName]*GitUser

	Save() error // TODO: this shouldn't be exposed
}

var persistentFlagHandle = func(c IssuectlConfig) error {
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

var inMemoryFlagHandle = func(_ IssuectlConfig) error { return nil }

func LoadConfig() IssuectlConfig {
	config := GetEmptyConfig()

	data, err := os.ReadFile(DefaultConfigFilePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := config.Save(); err != nil {
				return config
			}
			return config
		}
		return nil
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		return nil
	}

	return config
}

func GetEmptyConfig() IssuectlConfig {
	return &issuectlConfig{
		Repositories: map[RepoConfigName]*RepoConfig{},
		Issues:       map[IssueID]*IssueConfig{},
		Profiles:     map[ProfileName]*Profile{},
		Backends:     map[BackendConfigName]*BackendConfig{},
		GitUsers:     map[GitUserName]*GitUser{},
	}
}

func GetConfig(cn ProfileName, r map[RepoConfigName]*RepoConfig, b map[BackendConfigName]*BackendConfig, gu map[GitUserName]*GitUser, p map[ProfileName]*Profile) IssuectlConfig {
	return &issuectlConfig{
		CurrentProfile: cn,
		Repositories:   r,
		Profiles:       p,
		Backends:       b,
		GitUsers:       gu,
	}
}

func (ic *issuectlConfig) Save() error {
	switch ic._persistenceMode {
	case "persistent":
		return persistentFlagHandle(ic)
	case "memory":
		return inMemoryFlagHandle(ic)
	}

	return nil
}

func (ic *issuectlConfig) GetPersistent() IssuectlConfig {
	ic._persistenceMode = "persistent"
	return ic
}

func (ic *issuectlConfig) GetInMemory() IssuectlConfig {
	ic._persistenceMode = "memory"
	return ic
}

// Issues

func (ic *issuectlConfig) AddIssue(issueConfig *IssueConfig) error {
	ic.Issues[issueConfig.ID] = issueConfig
	return ic.Save()
}

func (ic *issuectlConfig) DeleteIssue(issueID IssueID) error {
	delete(ic.Issues, issueID)
	return ic.Save()
}

func (ic *issuectlConfig) GetIssue(issueID IssueID) (*IssueConfig, bool) {
	issue, ok := ic.Issues[issueID]
	return issue, ok
}

func (ic *issuectlConfig) GetIssues() map[IssueID]*IssueConfig {
	return ic.Issues
}

// Repositories

func (ic *issuectlConfig) GetRepository(name RepoConfigName) *RepoConfig {
	return ic.Repositories[name]
}

func (ic *issuectlConfig) AddRepository(repoConfig *RepoConfig) error {
	ic.Repositories[repoConfig.Name] = repoConfig
	return ic.Save()
}

func (ic *issuectlConfig) GetRepositories() map[RepoConfigName]*RepoConfig {
	return ic.Repositories
}

// Profiles

func (ic *issuectlConfig) GetProfile(profileName ProfileName) *Profile {
	return ic.Profiles[profileName]
}

func (ic *issuectlConfig) AddProfile(profile *Profile) error {
	ic.Profiles[profile.Name] = profile
	return ic.Save()
}

func (ic *issuectlConfig) DeleteProfile(profileName ProfileName) error {
	delete(ic.Profiles, profileName)
	return ic.Save()
}

func (ic *issuectlConfig) GetCurrentProfile() ProfileName {
	return ic.CurrentProfile
}

func (ic *issuectlConfig) UseProfile(profile ProfileName) error {
	ic.CurrentProfile = profile
	return ic.Save()
}

func (ic *issuectlConfig) GetProfiles() map[ProfileName]*Profile {
	return ic.Profiles
}

func (ic *issuectlConfig) UpdateProfile(profile *Profile) error {
	ic.Profiles[profile.Name] = profile
	return ic.Save()
}

// Backends
func (ic *issuectlConfig) GetBackend(backendName BackendConfigName) *BackendConfig {
	return ic.Backends[backendName]
}

func (ic *issuectlConfig) AddBackend(backend *BackendConfig) error {
	ic.Backends[backend.Name] = backend
	return ic.Save()
}

func (ic *issuectlConfig) DeleteBackend(backendName BackendConfigName) error {
	delete(ic.Backends, backendName)
	return ic.Save()
}

func (ic *issuectlConfig) GetBackends() map[BackendConfigName]*BackendConfig {
	return ic.Backends
}

// GitUsers
func (ic *issuectlConfig) GetGitUser(userName GitUserName) (*GitUser, bool) {
	gitUser, exists := ic.GitUsers[userName]
	return gitUser, exists
}

func (ic *issuectlConfig) AddGitUser(user *GitUser) error {
	ic.GitUsers[GitUserName(user.Name)] = user
	return ic.Save()
}

func (ic *issuectlConfig) DeleteGitUser(userName GitUserName) error {
	delete(ic.GitUsers, userName)
	return ic.Save()
}

func (ic *issuectlConfig) GetGitUsers() map[GitUserName]*GitUser {
	return ic.GitUsers
}
