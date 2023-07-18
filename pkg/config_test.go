package issuectl

// import (
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestIssuectlConfig_Save(t *testing.T) {
// 	config := &IssuectlConfig{
// 		urrentProfile: "testProfile",
// 	}

// 	err := config.Save()
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestLoadConfig(t *testing.T) {
// 	config := &IssuectlConfig{
// 		currentProfile: "testProfile",
// 	}

// 	err := config.Save()
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	loadedConfig := LoadConfig()
// 	if loadedConfig.currentProfile != config.currentProfile {
// 		t.Errorf("Expected CurrentProfile '%s', got '%s'", config.currentProfile, loadedConfig.currentProfile)
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestIssuectlConfig_AddIssue(t *testing.T) {
// 	config := &IssuectlConfig{}

// 	issue := &IssueConfig{
// 		ID: "testIssue",
// 	}

// 	err := config.AddIssue(issue)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if len(config.issues) != 1 {
// 		t.Errorf("Expected 1 issue, got %d", len(config.issues))
// 	}

// 	if config.issues[0].ID != issue.ID {
// 		t.Errorf("Expected issue ID '%s', got '%s'", issue.ID, config.issues[0].ID)
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestIssuectlConfig_DeleteIssue(t *testing.T) {
// 	config := &IssuectlConfig{}

// 	issue := &IssueConfig{
// 		ID: "testIssue",
// 	}

// 	config.AddIssue(issue) //nolint:errcheck

// 	assert.NoError(t, config.DeleteIssue(issue.ID))

// 	if len(config.issues) != 0 {
// 		t.Errorf("Expected 0 issues, got %d", len(config.issues))
// 	}

// 	// Test deleting a non-existing issue
// 	assert.Error(t, config.DeleteIssue("nonExistingIssue"))

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestIssuectlConfig_GetIssue(t *testing.T) {
// 	config := &IssuectlConfig{}

// 	issue := &IssueConfig{
// 		ID: "testIssue",
// 	}

// 	config.AddIssue(issue) //nolint:errcheck

// 	retrievedIssue := config.GetIssue(issue.ID)

// 	if retrievedIssue.ID != issue.ID { //nolint:staticcheck
// 		t.Errorf("Expected issue ID '%s', got '%s'", issue.ID, retrievedIssue.ID)
// 	}

// 	// Test getting a non-existing issue
// 	retrievedIssue = config.GetIssue("nonExistingIssue")
// 	if retrievedIssue != nil {
// 		t.Errorf("Expected nil, got issue")
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestIssuectlConfig_ListRepositories(t *testing.T) {
// 	config := &IssuectlConfig{}

// 	repo := &RepoConfig{
// 		Name: "testRepo",
// 	}

// 	config.AddRepository(repo) //nolint:errcheck

// 	err := config.ListRepositories()
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestIssuectlConfig_GetRepository(t *testing.T) {
// 	config := &IssuectlConfig{}

// 	repo := &RepoConfig{
// 		Name: "testRepo",
// 	}

// 	config.AddRepository(repo) //nolint:errcheck

// 	retrievedRepo := config.GetRepository(repo.Name)
// 	if retrievedRepo == nil { //nolint:staticcheck
// 		t.Errorf("Expected repository, got nil")
// 	}

// 	if retrievedRepo.Name != repo.Name { //nolint:staticcheck
// 		t.Errorf("Expected repository name '%s', got '%s'", repo.Name, retrievedRepo.Name)
// 	}

// 	// Test getting a non-existing repository
// 	retrievedRepo = config.GetRepository("nonExistingRepo")
// 	if retrievedRepo != nil {
// 		t.Errorf("Expected nil, got repository")
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }

// func TestIssuectlConfig_AddRepository(t *testing.T) {
// 	config := &IssuectlConfig{}

// 	repo := &RepoConfig{
// 		Name: "testRepo",
// 	}

// 	err := config.AddRepository(repo)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if len(config.repositories) != 1 {
// 		t.Errorf("Expected 1 repository, got %d", len(config.repositories))
// 	}

// 	if config.repositories[0].Name != repo.Name {
// 		t.Errorf("Expected repository name '%s', got '%s'", repo.Name, config.repositories[0].Name)
// 	}

// 	// Clean up
// 	os.Remove(DefaultConfigFilePath)
// }
