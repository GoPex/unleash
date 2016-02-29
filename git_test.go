package unleash_test

import (
	"os"
	"testing"

	// Unleash package to test
	"github.com/GoPex/unleash"
)

// Test the Clone function of the git helpers
func TestClone(t *testing.T) {
	defer os.RemoveAll(testDestinationPath + "_" + testRepositoryDefaultBranch)

	if _, err := unleash.Clone(testRepositoryUrl, testDestinationPath+"_"+testRepositoryDefaultBranch, testRepositoryDefaultBranch); err != nil {
		t.Error(err)
	}
}

// Test the Clone function of the git helpers with not the default branch
func TestCloneBranch(t *testing.T) {
	defer os.RemoveAll(testDestinationPath + "_" + testRepositoryNotDefaultBranch)

	if _, err := unleash.Clone(testRepositoryUrl, testDestinationPath+"_"+testRepositoryNotDefaultBranch, testRepositoryNotDefaultBranch); err != nil {
		t.Error(err)
	}
}

// Test the Clone function of the git helpers
func TestClonePrivate(t *testing.T) {
	defer os.RemoveAll(testDestinationPath + "_private")

	if _, err := unleash.Clone(testPrivateRepositoryUrl, testDestinationPath+"_private", testRepositoryDefaultBranch); err != nil {
		t.Error(err)
	}
}
