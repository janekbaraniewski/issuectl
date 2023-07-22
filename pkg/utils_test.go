package issuectl

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestCloneRepo tests the cloneRepo function.
func TestCloneRepo(t *testing.T) {
	// Mocking the RepoConfig and GitUser
	repo := &RepoConfig{Name: "testRepo", RepoURL: "https://github.com/janekbaraniewski/issuectl.git"}
	gitUser := &GitUser{Name: "testUser", Email: "test@example.com", SSHKey: "/dev/null"}

	// Create a temporary directory to clone the repo
	dir, err := ioutil.TempDir("", "testDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(dir) // clean up

	// Call the cloneRepo function
	_, err = cloneRepo(repo, dir, gitUser)
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
	dir, err := ioutil.TempDir("", "testDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(dir) // clean up

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
	parentDir, err := ioutil.TempDir("", "parentDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(parentDir) // clean up

	// Call the createDirectory function
	_, err = createDirectory(parentDir, "testDir")
	if err != nil {
		t.Fatalf("createDirectory() failed: %s", err)
	}
}
