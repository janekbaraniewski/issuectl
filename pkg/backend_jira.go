package issuectl

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

type Jira struct {
	baseURL  string
	email    string
	apiToken string
	client   *jira.Client
}

const (
	ToDo       = "To Do"
	InProgress = "In Progress"
	Done       = "Done"
)

func NewJiraClient(email, apiToken, baseURL string) *Jira {
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiToken,
	}

	client, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		panic(err)
	}

	return &Jira{client: client, baseURL: baseURL, email: email, apiToken: apiToken}
}

func (j *Jira) GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error) {
	issue, _, err := j.client.Issue.Get(string(issueID), nil)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (j *Jira) StartIssue(owner string, repo RepoConfigName, issueID IssueID) error {
	return j.moveIssueToState(issueID, InProgress)
}

func (j *Jira) CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error {
	return j.moveIssueToState(issueID, Done)
}

func (j *Jira) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error {
	link := jira.IssueLink{
		Type: jira.IssueLinkType{
			Name:    "Cloners",
			Inward:  "clones",
			Outward: "is cloned by",
		},
		InwardIssue:  &jira.Issue{Key: string(issueID)},
		OutwardIssue: &jira.Issue{Key: pullRequestID},
	}
	_, err := j.client.Issue.AddLink(&link)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jira) moveIssueToState(issueID IssueID, desiredState string) error {
	issue, _, err := j.client.Issue.Get(string(issueID), nil)
	if err != nil {
		return err
	}

	// If the issue is already in the desired state, return
	if issue.Fields.Status.Name == desiredState {
		return nil
	}

	transitions, _, err := j.client.Issue.GetTransitions(string(issueID))
	if err != nil {
		return err
	}

	// Find the transition to the desired state
	var transitionID string
	for _, transition := range transitions {
		if transition.Name == desiredState {
			transitionID = transition.ID
			break
		}
	}

	if transitionID == "" {
		return fmt.Errorf("unable to find '%s' transition", desiredState)
	}

	_, err = j.client.Issue.DoTransition(string(issueID), transitionID)
	if err != nil {
		return err
	}

	return nil
}
