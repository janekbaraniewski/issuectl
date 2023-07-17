package issuectl

import "os/exec"

func cloneRepo(url, dir string) error {
	cmd := exec.Command("git", "clone", url, dir)
	return cmd.Run()
}

func createBranch(dir, branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dir
	return cmd.Run()
}
