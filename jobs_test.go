package unleash_test

import (
	"testing"

	// Unleash package to test
	"github.com/GoPex/unleash"
)

// Test the BuildAndPushFromRepository job
func TestBuildAndPushFromRepository(t *testing.T) {
	defer dockerClient.RemoveImage(testImageRepository+":latest", true)

	if err := unleash.BuildAndPushFromRepository(testRepositoryUrl, testRepositoryFullName, testRepositoryDefaultBranch, testRepositoryCommitId); err != nil {
		t.Error(err)
	}
}
