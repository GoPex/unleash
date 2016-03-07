package jobs_test

import (
	"testing"

	"github.com/GoPex/dockerclient"
	"github.com/GoPex/unleash/jobs"
	"github.com/GoPex/unleash/tests"
)

var (
	dockerClient, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
)

// Test the BuildAndPushFromRepository job
func TestBuildAndPushFromRepository(t *testing.T) {
	defer dockerClient.RemoveImage(tests.TestImageRepository+":latest", true)

	if err := jobs.BuildAndPushFromRepository(tests.GithubArchiveURL(tests.TestRepositoryName, tests.TestRepositoryDefaultBranch), tests.TestRepositoryFullName, tests.TestRepositoryDefaultBranch, tests.TestRepositoryCommitID); err != nil {
		t.Error(err)
	}
}
