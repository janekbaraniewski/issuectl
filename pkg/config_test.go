package issuectl_test

import (
	"testing"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGetRepository(t *testing.T) {
	config := issuectl.LoadConfig()

	repo := config.GetRepository("multi-cloud")
	assert.NotNil(t, repo)

	repo = config.GetRepository("some-name-that-doesnt-exist")
	assert.Nil(t, repo)
}
