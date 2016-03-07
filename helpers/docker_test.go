package helpers_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoPex/dockerclient"
	"github.com/Rolinh/targo"

	"github.com/GoPex/unleash/helpers"
	"github.com/GoPex/unleash/tests"
)

var (
	dockerClient, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

	testRepositoryTarPath                  = filepath.Join(tests.DataDirectory, "unleash_test_repository.tar")
	testUnknowInstructionRepositoryTarPath = filepath.Join(tests.DataDirectory, "unleash_test_repository_unknown_instruction.tar")
	testNonZeroCodeRepositoryTarPath       = filepath.Join(tests.DataDirectory, "unleash_test_repository_non-zero_code.tar")

	testRepositoryExtracted = filepath.Join(tests.WorkingDirectory, "unleash_test_repository_extracted")
)

// Test the BuildFromTar function of the docker helpers
func TestBuildFromTar(t *testing.T) {
	defer dockerClient.RemoveImage(tests.TestImageRepository+":fromTar", true)

	id, err := helpers.BuildFromTar(testRepositoryTarPath, tests.TestImageRepository+":fromTar", tests.ContextLogger)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("No id was returned !")
	}
}

// Test unknown instruction error handling of the BuildFromTar function of the docker helpers
func TestBuildFromTarUnknowInstruction(t *testing.T) {
	id, err := helpers.BuildFromTar(testUnknowInstructionRepositoryTarPath, tests.TestImageRepository+":buildError", tests.ContextLogger)
	if err == nil {
		t.Error("No error returned ! It should return an error as the build should fail !")
	}
	if id != "" {
		t.Error("An id was returned and it shouldn't as the build should fail !")
	}
}

// Test non zero code error handling of the BuildFromTar function of the docker helpers
func TestBuildFromTarNonZeroCode(t *testing.T) {
	id, err := helpers.BuildFromTar(testNonZeroCodeRepositoryTarPath, tests.TestImageRepository+":buildError", tests.ContextLogger)
	if err == nil {
		t.Error("No error returned ! It should return an error as the build should fail !")
	}
	if id != "" {
		t.Error("An id was returned and it shouldn't as the build should fail !")
	}
}

// Test the BuildFromDirectory function of the docker helpers
func TestBuildFromDirectory(t *testing.T) {
	defer dockerClient.RemoveImage(tests.TestImageRepository+":fromDirectory", true)
	defer os.RemoveAll(testRepositoryExtracted)

	targo.Extract(testRepositoryExtracted, testRepositoryTarPath)
	id, err := helpers.BuildFromDirectory(testRepositoryExtracted, tests.TestImageRepository+":fromDirectory", tests.ContextLogger)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("No id was returned !")
	}
}

// Test the PushImage function of the docker helpers
func TestPushImage(t *testing.T) {
	testImageFullRepository := tests.TestDockerRegistryURL + "/" + tests.TestImageRepository + ":testPush"
	defer dockerClient.RemoveImage(testImageFullRepository, true)

	// It's not ideal to rely on our function for this test but its simpler for now
	_, err := helpers.BuildFromTar(testRepositoryTarPath, testImageFullRepository, tests.ContextLogger)
	if err != nil {
		t.Error(err)
	} else {
		if err := helpers.PushImage(testImageFullRepository); err != nil {
			t.Error(err)
		}
	}
}

// Test the Ping function of the docker helpers
func TestPing(t *testing.T) {
	pong, err := helpers.Ping()
	if err != nil {
		t.Error(err)
	}
	if pong != "OK" {
		t.Error(errors.New("Pong response is not OK but " + pong))
	}
}

// Test the Version function of the Docker helpers
func TestVersion(t *testing.T) {
	version, err := helpers.Version()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error(errors.New("Version is empty !"))
	}
}
