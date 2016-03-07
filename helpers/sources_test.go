package helpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/GoPex/unleash/helpers"
	"github.com/GoPex/unleash/tests"
)

var (
	testDestinationPath = filepath.Join(tests.WorkingDirectory, "unleash_test_repository")
)

// TestExtractRepository tests the ExtractRepository function
func TestExtractRepository(t *testing.T) {
	workingDirectory := testDestinationPath + "_" + tests.TestRepositoryDefaultBranch
	os.MkdirAll(workingDirectory, 0600)
	defer os.RemoveAll(workingDirectory)

	if err := helpers.ExtractRepository(tests.GithubArchiveURL(tests.TestRepositoryName, tests.TestRepositoryDefaultBranch), workingDirectory); err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(workingDirectory + "_bitbucket")
	if err := helpers.ExtractRepository(tests.BitbucketArchiveURL(tests.TestRepositoryName, tests.TestRepositoryDefaultBranch), workingDirectory+"_bitbucket"); err != nil {
		t.Error(err)
	}
}

// TestExtractRepository_private tests the ExtractRepository function with a private repository
func TestExtractRepository_private(t *testing.T) {
	workingDirectory := testDestinationPath + "_" + tests.TestRepositoryDefaultBranch + "_private"
	defer os.RemoveAll(workingDirectory)
	os.MkdirAll(workingDirectory, 0600)

	if err := helpers.ExtractRepository(tests.BitbucketArchiveURL(tests.TestPrivateRepositoryName, tests.TestRepositoryNotDefaultBranch), workingDirectory); err != nil {
		t.Error(err)
	}
}
