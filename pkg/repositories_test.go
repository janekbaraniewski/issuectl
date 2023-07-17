package issuectl_test

import (
	"testing"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGetRepository(t *testing.T) {
	config := issuectl.GetRepository("multi-cloud")
	assert.NotNil(t, config)

	config = issuectl.GetRepository("some-name-that-doesnt-exist")
	assert.Nil(t, config)
}
