package issuectl

import (
	"os"
	"os/exec"
	"path/filepath"
)

func cloneRepo(repo *RepoConfig, dir string) (string, error) {
	repoDir := filepath.Join(dir, string(repo.Name))
	cmd := exec.Command("git", "clone", string(repo.RepoUrl), repoDir)
	return repoDir, cmd.Run()
}

func createBranch(dir, branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dir
	return cmd.Run()
}

func createDirectory(parent, dirName string) (string, error) {
	dirPath := filepath.Join(parent, dirName)
	if err := os.Mkdir(dirPath, 0755); err != nil {
		return "", err
	}
	return dirPath, nil
}
