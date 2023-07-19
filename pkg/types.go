package issuectl

// BackendType is a name of issue backend
type BackendType string

// BackendGithub is a BackendType for GitHub
const (
	BackendGithub BackendType = "github"
	BackendGitLab BackendType = "gitlab"
)

// BackendConfigName is a name of instance of BackendConfig
type BackendConfigName string

// BackendConfig stores configuration for given BackendType
type BackendConfig struct {
	// Name of BackendConfig instance
	Name BackendConfigName `json:"name"`

	// BackendType of this BackendConfig
	Type BackendType `json:"backendType"`

	// Path to SSL key to use for git
	SSLKeyPath string `json:"auth"`
}

// RepoURL is a string with URL to git repo for cloning
type RepoURL string

// RepoConfigName is a name of git repo
type RepoConfigName string

// RepoConfig stores configuration for given git repo
type RepoConfig struct {
	// Name of this repo
	Name RepoConfigName `json:"name"`

	// Repo owner
	Owner string `json:"owner"`

	// URL to this repo
	RepoURL RepoURL `json:"url"`
}

// IssueID is a unique ID of issue in IssueBackend
type IssueID string

// IssueConfig stores configuration for single issue
type IssueConfig struct {
	Name         string            `json:"name"`
	ID           IssueID           `json:"id"`
	BackendName  BackendConfigName `json:"backendName"`
	BranchName   string            `json:"branchName"`
	RepoName     RepoConfigName    `json:"repoName"`
	Repositories []RepoConfigName  `json:"repositories"`
	Dir          string            `json:"dir"`
	Profile      ProfileName       `json:"profile"`
}
