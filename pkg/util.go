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
	Log.V(3).Infof("git clone %v %v", repo.RepoURL, repoDir)
	cmd := exec.Command("git", "clone", string(repo.RepoURL), repoDir)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	if err := setRepoIdentity(repoDir, gitUser.Name, gitUser.Email, gitUser.SSHKey); err != nil {
		return "", err
	}

	return repoDir, nil
}

// createBranch takes a directory, a branch name, and a GitUser object as arguments.
// It creates a new git branch with the specified name in the specified directory.
// It returns any error encountered during the branch creation process.
func createBranch(dir, branchName string, gitUser *GitUser) error {
	if err := setRepoIdentity(dir, gitUser.Name, gitUser.Email, gitUser.SSHKey); err != nil {
		return err
	}

	exists, err := branchExists(dir, branchName)
	if err != nil {
		return err
	}

	if exists {
		Log.V(3).Infof("git checkout %v", branchName)
		cmd := exec.Command("git", "checkout", branchName)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		Log.V(3).Infof("git checkout -b %v", branchName)
		cmd := exec.Command("git", "checkout", "-b", branchName)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return err
		}

		Log.V(3).Infof("git push --set-upstream origin %v", branchName)
		cmd = exec.Command("git", "push", "--set-upstream", "origin", branchName)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// branchExists checks if a branch exists in the repository located at dir.
func branchExists(dir, branchName string) (bool, error) {
	cmd := exec.Command("git", "branch", "--list", branchName)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

// setRepoIdentity sets local git config username, email and ssh command.
func setRepoIdentity(dir string, username GitUserName, email, sshKeyPath string) error {
	// Set local git config user.name
	Log.V(3).Infof("git config user.name %v", username)
	cmd := exec.Command("git", "config", "user.name", string(username))
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set local git config user.email
	Log.V(3).Infof("git config user.email %v", email)
	cmd = exec.Command("git", "config", "user.email", email)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set local git config core.sshCommand
	sshCommand := fmt.Sprintf("ssh -i %s -F /dev/null", sshKeyPath)
	Log.V(3).Infof("git config core.sshCommand %v", sshCommand)
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
