package issuectl

import (
	"os"
	"testing"
)

// TestCloneRepo tests the cloneRepo function.
func TestCloneRepo(t *testing.T) {
	// Mocking the RepoConfig and GitUser
	repo := &RepoConfig{Name: "testRepo", RepoURL: "https://github.com/janekbaraniewski/issuectl.git"}
	gitUser := &GitUser{Name: "testUser", Email: "test@example.com", SSHKey: "/dev/null"}

	// Create a temporary directory to clone the repo
	dir := os.TempDir()     //nolint:gosimple
	defer os.RemoveAll(dir) // clean up

	// Call the cloneRepo function
	_, err := cloneRepo(repo, dir, gitUser)
	if err != nil {
		t.Fatalf("cloneRepo() failed: %s", err)
	}
}

// TestCreateBranch tests the createBranch function.
func TestCreateBranch(t *testing.T) {
	t.Skip("FIXME")
	// Mocks
	repo := &RepoConfig{Name: "testRepo", RepoURL: "https://github.com/janekbaraniewski/issuectl.git"}
	gitUser := &GitUser{Name: "testUser", Email: "test@example.com", SSHKey: "/dev/null"}

	// Create a temporary directory to clone the repo
	dir := os.TempDir()     //nolint:gosimple
	defer os.RemoveAll(dir) // nolint

	// Call the cloneRepo function
	repoDir, err := cloneRepo(repo, dir, gitUser)
	if err != nil {
		t.Fatalf("cloneRepo() failed: %s", err)
	}

	// Call the createBranch function
	if err := createBranch(repoDir, "testBranch", gitUser); err != nil {
		t.Fatalf("createBranch() failed: %s", err)
	}
}

// TestCreateDirectory tests the createDirectory function.
func TestCreateDirectory(t *testing.T) {
	// Create a temporary directory
	parentDir := os.TempDir()     //nolint:gosimple
	defer os.RemoveAll(parentDir) // nolint

	// Call the createDirectory function
	_, err := createDirectory(parentDir, "testDir")
	if err != nil {
		t.Fatalf("createDirectory() failed: %s", err)
	}
}
