package issuectl

var multiCloudRepo = &RepoConfig{
	Name:    "multi-cloud",
	RepoUrl: "git@github.com:elotl/multi-cloud.git",
}

var WorkDir string = "/Users/janbaraniewski/Workspace/priv/issuectl/testWorkdir" // FIXME

func StartWorkingOnIssue(issueID IssueID) error {
	Log.Infof("Starting work on issue %v ...", issueID)
	Log.V(2).Infof("Creating issue work dir")
	issueDirPath, err := createDirectory(WorkDir, string(issueID))
	if err != nil {
		return err
	}
	Log.V(2).Infof("Cloning repo")
	repoDirPath, err := cloneRepo(multiCloudRepo, issueDirPath)
	if err != nil {
		return err
	}

	Log.V(2).Infof("Creating branch")
	if err = createBranch(repoDirPath, string(issueID)); err != nil {
		return err
	}

	Log.Infof("Started working on issue %v", issueID)

	return nil
}
