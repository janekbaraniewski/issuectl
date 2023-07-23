package issuectl

// BackendType is a name of issue backend
type BackendType string

// BackendGithub is a BackendType for GitHub
const (
	BackendGithub BackendType = "github"
	BackendGitLab BackendType = "gitlab"
	BackendJira   BackendType = "jira"
)

// BackendConfigName is a name of instance of BackendConfig
type BackendConfigName string

type GitHubConfig struct {
	Host     string `yaml:"host,omitempty"`
	Token    string `yaml:"token,omitempty"`
	Username string `yaml:"username,omitempty"`
}

type GitLabConfig struct {
	Host   string `yaml:"host,omitempty"`
	Token  string `yaml:"token,omitempty"`
	UserID int    `yaml:"userID,omitempty"`
}

type JiraConfig struct {
	Host     string `yaml:"host,omitempty"`
	Token    string `yaml:"token,omitempty"`
	Username string `yaml:"username,omitempty"`
}

// BackendConfig stores configuration for given BackendType
type BackendConfig struct {
	// Name of BackendConfig instance
	Name BackendConfigName `yaml:"name"`

	// BackendType of this BackendConfig
	Type BackendType `yaml:"backendType"`

	GitHub *GitHubConfig `yaml:"github,omitempty"`
	GitLab *GitLabConfig `yaml:"gitlab,omitempty"`
	Jira   *JiraConfig   `yaml:"jira,omitempty"`
}

type GitUserName string

// GitUser holds config for git user
type GitUser struct {
	Name   GitUserName
	Email  string
	SSHKey string
}

// RepoURL is a string with URL to git repo for cloning
type RepoURL string

// RepoConfigName is a name of git repo
type RepoConfigName string

// RepoConfig stores configuration for given git repo
type RepoConfig struct {
	// Name of this repo
	Name RepoConfigName `yaml:"name"`

	// Repo owner
	Owner string `yaml:"owner"`

	// URL to this repo
	RepoURL RepoURL `yaml:"url"`
}

// IssueID is a unique ID of issue in IssueBackend
type IssueID string

// IssueConfig stores configuration for single issue
type IssueConfig struct {
	Name         string            `yaml:"name"`
	ID           IssueID           `yaml:"id"`
	BackendName  BackendConfigName `yaml:"backendName"`
	BranchName   string            `yaml:"branchName"`
	Repositories []RepoConfigName  `yaml:"repositories"`
	Dir          string            `yaml:"dir"`
	Profile      ProfileName       `yaml:"profile"`
}
