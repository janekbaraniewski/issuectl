package issuectl

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type GitLab struct {
	client *gitlab.Client
}

func NewGitLabClient(token, baseURL string) *GitLab {
	Log.Infof("NOT TESTED")
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		Log.Infof("failed to create GitLab client: %v", err)
		return nil
	}

	return &GitLab{client: client}
}

func (g *GitLab) GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error) {
	Log.Infof("NOT IMPLEMENTED")
	return nil, nil
}

func (g *GitLab) CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return err
	}

	issueOpt := &gitlab.UpdateIssueOptions{
		StateEvent: gitlab.String("close"),
	}

	_, _, err = g.client.Issues.UpdateIssue(fmt.Sprintf("%s/%s", owner, repo), issueNumber, issueOpt)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitLab) OpenPullRequest(owner string, repo RepoConfigName, title, body, baseBranch, headBranch string) error {
	pullReqOpt := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String(title),
		Description:  gitlab.String(body),
		SourceBranch: gitlab.String(headBranch),
		TargetBranch: gitlab.String(baseBranch),
	}

	_, _, err := g.client.MergeRequests.CreateMergeRequest(fmt.Sprintf("%s/%s", owner, repo), pullReqOpt)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitLab) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return err
	}

	issueRef := fmt.Sprintf("%s/%s#%d", owner, repo, issueNumber)

	pullReqOpt := &gitlab.UpdateMergeRequestOptions{
		Description: gitlab.String(fmt.Sprintf("Closes %s", issueRef)),
	}

	pullRequestNumber, err := getIssueNumberFromString(IssueID(pullRequestID))
	if err != nil {
		return err
	}

	_, _, err = g.client.MergeRequests.UpdateMergeRequest(fmt.Sprintf("%s/%s", owner, repo), pullRequestNumber, pullReqOpt)
	if err != nil {
		return err
	}

	return nil
}
