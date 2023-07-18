package issuectl

import (
	"os"
	"os/exec"
	"path/filepath"
)

// cloneRepo takes a RepoConfig object and a directory name as arguments.
// It clones the repository URL from the RepoConfig into the specified directory,
// and returns the path of the new repository directory and any error encountered.
func cloneRepo(repo *RepoConfig, dir string) (string, error) {
	repoDir := filepath.Join(dir, string(repo.Name))
	cmd := exec.Command("git", "clone", string(repo.RepoURL), repoDir)
	return repoDir, cmd.Run()
}

// createBranch takes a directory and a branch name as arguments.
// It creates a new git branch with the specified name in the specified directory.
// It returns any error encountered during the branch creation process.
func createBranch(dir, branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("git", "push", "--set-upstream", "origin", branchName)
	cmd.Dir = dir
	return cmd.Run()
}

// createDirectory takes a parent directory and a directory name as arguments.
// It creates a new directory with the specified name inside the parent directory.
// It returns the path of the new directory and any error encountered.
func createDirectory(parent, dirName string) (string, error) {
	dirPath := filepath.Join(parent, dirName)
	if err := os.Mkdir(dirPath, 0755); err != nil {
		return "", err
	}
	return dirPath, nil
}
