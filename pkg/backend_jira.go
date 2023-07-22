package issuectl

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
)

type Jira struct {
	client *jira.Client
}

func NewJiraBackend(username, token, host string) *Jira {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), host)
	if err != nil {
		return nil
	}

	return &Jira{
		client: client,
	}
}

func (j *Jira) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error {
	// TODO: add comment with link to PR?
	return nil
}

func (j *Jira) CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error {
	issue, _, err := j.client.Issue.Get(string(issueID), nil)
	if err != nil {
		return err
	}
	currentStatus := issue.Fields.Status.Name
	fmt.Printf("Current status: %s\n", currentStatus)

	var transitionID string
	possibleTransitions, _, err := j.client.Issue.GetTransitions(string(issueID))
	if err != nil {
		return err
	}
	for _, v := range possibleTransitions {
		if v.Name == "In Progress" {
			transitionID = v.ID
			break
		}
	}

	if _, err := j.client.Issue.DoTransition(string(issueID), transitionID); err != nil {
		return err
	}
	issue, _, err = j.client.Issue.Get(string(issueID), nil)
	if err != nil {
		return err
	}
	fmt.Printf("Status after transition: %+v\n", issue.Fields.Status.Name)
	return nil
}

func (j *Jira) GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error) {
	issue, _, err := j.client.Issue.Get(string(issueID), nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
func (j *Jira) IssueExists(owner string, repo RepoConfigName, issueID IssueID) (bool, error) {
	issue, err := j.GetIssue(owner, repo, issueID)
	return issue == nil, err
}
