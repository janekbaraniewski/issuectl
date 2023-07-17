package issuectl_test

import (
	"testing"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGetRepository(t *testing.T) {
<<<<<<< HEAD
	config := issuectl.LoadConfig()

	repo := config.GetRepository("multi-cloud")
	assert.NotNil(t, repo)

	repo = config.GetRepository("some-name-that-doesnt-exist")
	assert.Nil(t, repo)
=======
	config := issuectl.GetRepository("multi-cloud")
	assert.NotNil(t, config)

	config = issuectl.GetRepository("some-name-that-doesnt-exist")
	assert.Nil(t, config)
>>>>>>> 99b6603 (Keep all changes to config object in one place)
}
