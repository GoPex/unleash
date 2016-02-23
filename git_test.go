package unleash_test

import (
    "os"
    "testing"

    // Unleash package to test
    "bitbucket.org/gopex/unleash"
)

// Test the Clone function of the git helpers
func TestClone(t *testing.T) {
    defer os.RemoveAll(testDestinationPath)

    if err := unleash.Clone(testRepositoryUrl, testDestinationPath, "master"); err != nil {
        t.Error(err)
    }
}
