package issuectl

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGitHub_IssueExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected method 'GET', got '%s'", req.Method)
		}

		if req.URL.String() != "/repos/owner/repo/issues/1" {
			t.Errorf("Expected URL '/repos/owner/repo/issues/1', got '%s'", req.URL.String())
		}

		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	g := &GitHub{}

	exists, err := g.IssueExists(server.URL, "owner", "repo", "1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !exists {
		t.Errorf("Expected issue to exist, got %v", exists)
	}
}

func TestGitHub_LinkIssueToRepo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Errorf("Expected method 'POST', got '%s'", req.Method)
		}

		if req.URL.String() != "/repos/owner/repo/issues/1/timeline" {
			t.Errorf("Expected URL '/repos/owner/repo/issues/1/timeline', got '%s'", req.URL.String())
		}

		if req.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Expected header 'Authorization: Bearer token', got 'Authorization: %s'", req.Header.Get("Authorization"))
		}

		if req.Header.Get("Accept") != "application/vnd.github.starfire-preview+json" {
			t.Errorf("Expected header 'Accept: application/vnd.github.starfire-preview+json', got 'Accept: %s'", req.Header.Get("Accept"))
		}

		body := make([]byte, req.ContentLength)
		req.Body.Read(body) //nolint:errcheck
		if !strings.Contains(string(body), `"issue_number":"1"`) {
			t.Errorf("Expected body to contain '\"issue_number\":\"1\"', got '%s'", string(body))
		}

		rw.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	g := &GitHub{}

	err := g.LinkIssueToRepo(server.URL, "owner", "repo", "1", "1", "token")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGitHub_CloseIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			t.Errorf("Expected method 'PATCH', got '%s'", req.Method)
		}

		if req.URL.String() != "/repos/owner/repo/issues/1" {
			t.Errorf("Expected URL '/repos/owner/repo/issues/1', got '%s'", req.URL.String())
		}

		if req.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Expected header 'Authorization: Bearer token', got 'Authorization: %s'", req.Header.Get("Authorization"))
		}

		if req.Header.Get("Accept") != "application/vnd.github.v3+json" {
			t.Errorf("Expected header 'Accept: application/vnd.github.v3+json', got 'Accept: %s'", req.Header.Get("Accept"))
		}

		body := make([]byte, req.ContentLength)
		req.Body.Read(body) //nolint:errcheck
		if !strings.Contains(string(body), `"state":"closed"`) {
			t.Errorf("Expected body to contain '\"state\":\"closed\"', got '%s'", string(body))
		}

		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	g := &GitHub{}

	err := g.CloseIssue(server.URL, "owner", "repo", "1", "token")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
