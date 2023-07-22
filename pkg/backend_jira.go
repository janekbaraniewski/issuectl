package issuectl

import (
	jira "github.com/andygrunwald/go-jira"
)

type Jira struct {
	client *jira.Client
}

func NewJiraBackend(username, token string) *Jira {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), "https://my.jira.com")
	if err != nil {
		return nil
	}

	return &Jira{
		client: client,
	}
}

func (j *Jira) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error
func (j *Jira) CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error
func (j *Jira) GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error) {
	issue, _, err := j.client.Issue.Get("MESOS-3325", nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
func (j *Jira) IssueExists(owner string, repo RepoConfigName, issueID IssueID) (bool, error)
