package unleash_test

import (
	"testing"

	"github.com/GoPex/unleash"
)

// Test the BuildAndPushFromRepository job
func TestBuildAndPushFromRepository(t *testing.T) {
	defer dockerClient.RemoveImage(testImageRepository+":latest", true)

	if err := unleash.BuildAndPushFromRepository(githubArchiveUrl(testRepositoryName, testRepositoryDefaultBranch), testRepositoryFullName, testRepositoryDefaultBranch, testRepositoryCommitId); err != nil {
		t.Error(err)
	}
}
