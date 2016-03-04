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

// TestExtractRepository_bitbucket tests the ExtractRepository function
func TestExtractRepository_bitbucket(t *testing.T) {
	workingDirectory := testDestinationPath + "_" + testRepositoryDefaultBranch
	defer os.RemoveAll(workingDirectory)
	os.MkdirAll(workingDirectory, 0600)

	if err := unleash.ExtractRepository(bitbucketArchiveUrl(testRepositoryName, testRepositoryDefaultBranch), workingDirectory); err != nil {
		t.Error(err)
	}
}

// TestExtractRepository_private tests the ExtractRepository function with a private repository
func TestExtractRepository_private(t *testing.T) {
	workingDirectory := testDestinationPath + "_" + testRepositoryDefaultBranch + "_private"
	defer os.RemoveAll(workingDirectory)
	os.MkdirAll(workingDirectory, 0600)

	if err := unleash.ExtractRepository(bitbucketArchiveUrl(testPrivateRepositoryName, testRepositoryNotDefaultBranch), workingDirectory); err != nil {
		t.Error(err)
	}
}
