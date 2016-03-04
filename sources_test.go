package unleash_test

import (
	"os"
	"testing"

	"github.com/GoPex/unleash"
)

// TestExtractRepository tests the ExtractRepository function
func TestExtractRepository(t *testing.T) {
	workingDirectory := testDestinationPath + "_" + testRepositoryDefaultBranch
	defer os.RemoveAll(workingDirectory)
	os.MkdirAll(workingDirectory, 0600)

	if err := unleash.ExtractRepository(githubArchiveUrl(testRepositoryName, testRepositoryDefaultBranch), workingDirectory); err != nil {
		t.Error(err)
	}
}
