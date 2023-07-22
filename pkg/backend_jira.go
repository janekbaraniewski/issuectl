package issuectl

type Jira struct{}

func NewJiraBackend() *Jira {
	return &Jira{}
}

func (j *Jira) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error
func (j *Jira) CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error
func (j *Jira) GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error)
func (j *Jira) IssueExists(owner string, repo RepoConfigName, issueID IssueID) (bool, error)
