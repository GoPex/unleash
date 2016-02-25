package unleash_test

import (
    "os"
    "testing"

    // Unleash package to test
    "bitbucket.org/gopex/unleash"
)

// Test the Clone function of the git helpers
func TestClone(t *testing.T) {
    defer os.RemoveAll(testDestinationPath + "_" + testRepositoryDefaultBranch)

    if err := unleash.Clone(testRepositoryUrl, testDestinationPath + "_" + testRepositoryDefaultBranch, testRepositoryDefaultBranch); err != nil {
        t.Error(err)
    }
}

// Test the Clone function of the git helpers with not the default branch
func TestCloneBranch(t *testing.T) {
    defer os.RemoveAll(testDestinationPath + "_" + testRepositoryNotDefaultBranch)

    if err := unleash.Clone(testRepositoryUrl, testDestinationPath + "_" + testRepositoryNotDefaultBranch, testRepositoryNotDefaultBranch); err != nil {
        t.Error(err)
    }
}
