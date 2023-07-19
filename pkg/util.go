package issuectl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// cloneRepo takes a RepoConfig object, a directory name, and a GitUser object as arguments.
// It clones the repository URL from the RepoConfig into the specified directory,
// and returns the path of the new repository directory and any error encountered.
func cloneRepo(repo *RepoConfig, dir string, gitUser *GitUser) (string, error) {
	repoDir := filepath.Join(dir, string(repo.Name))
	cmd := exec.Command("git", "clone", string(repo.RepoURL), repoDir)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	if err := setRepoIdentity(repoDir, gitUser.GitUserName, gitUser.Email, gitUser.SSHKey); err != nil {
		return "", err
	}

	return repoDir, nil
}

// createBranch takes a directory, a branch name, and a GitUser object as arguments.
// It creates a new git branch with the specified name in the specified directory.
// It returns any error encountered during the branch creation process.
func createBranch(dir, branchName string, gitUser *GitUser) error {
	if err := setRepoIdentity(dir, gitUser.GitUserName, gitUser.Email, gitUser.SSHKey); err != nil {
		return err
	}

	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "push", "--set-upstream", "origin", branchName)
	cmd.Dir = dir
	return cmd.Run()
}

// setRepoIdentity sets local git config username, email and ssh command.
func setRepoIdentity(dir, username, email, sshKeyPath string) error {
	// Set local git config user.name
	cmd := exec.Command("git", "config", "user.name", username)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set local git config user.email
	cmd = exec.Command("git", "config", "user.email", email)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set local git config core.sshCommand
	sshCommand := fmt.Sprintf("ssh -i %s -F /dev/null", sshKeyPath)
	cmd = exec.Command("git", "config", "core.sshCommand", sshCommand)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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
