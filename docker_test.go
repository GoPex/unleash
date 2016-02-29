package unleash_test

import (
	"errors"
	"os"
	"testing"

	// go client for docker
	"github.com/GoPex/dockerclient"

	// Tar utilities for go
	"github.com/Rolinh/targo"

	// Unleash package to test
	"github.com/GoPex/unleash"
)

var (
	dockerClient, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
)

// Test the BuildFromTar function of the docker helpers
func TestBuildFromTar(t *testing.T) {
	defer dockerClient.RemoveImage(testImageRepository+":fromTar", true)

	id, err := unleash.BuildFromTar(testRepositoryTarPath, testImageRepository+":fromTar", contextLogger)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("No id was returned !")
	}
}

// Test unknown instruction error handling of the BuildFromTar function of the docker helpers
func TestBuildFromTarUnknowInstruction(t *testing.T) {
	id, err := unleash.BuildFromTar(testUnknowInstructionRepositoryTarPath, testImageRepository+":buildError", contextLogger)
	if err == nil {
		t.Error("No error returned ! It should return an error as the build should fail !")
	}
	if id != "" {
		t.Error("An id was returned and it shouldn't as the build should fail !")
	}
}

// Test non zero code error handling of the BuildFromTar function of the docker helpers
func TestBuildFromTarNonZeroCode(t *testing.T) {
	id, err := unleash.BuildFromTar(testNonZeroCodeRepositoryTarPath, testImageRepository+":buildError", contextLogger)
	if err == nil {
		t.Error("No error returned ! It should return an error as the build should fail !")
	}
	if id != "" {
		t.Error("An id was returned and it shouldn't as the build should fail !")
	}
}

// Test the BuildFromDirectory function of the docker helpers
func TestBuildFromDirectory(t *testing.T) {
	defer dockerClient.RemoveImage(testImageRepository+":fromDirectory", true)
	defer os.RemoveAll(testRepositoryExtracted)

	targo.Extract(testRepositoryExtracted, testRepositoryTarPath)
	id, err := unleash.BuildFromDirectory(testRepositoryExtracted, testImageRepository+":fromDirectory", contextLogger)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("No id was returned !")
	}
}

// Test the PushImage function of the docker helpers
func TestPushImage(t *testing.T) {
	defer dockerClient.RemoveImage(testImageRepository+":testPush", true)

	// It's not ideal to rely on our function for this test but its simpler for now
	_, err := unleash.BuildFromTar(testRepositoryTarPath, testImageRepository+":testPush", contextLogger)
	if err != nil {
		t.Error(err)
	} else {
		if err := unleash.PushImage(testImageRepository + ":testPush"); err != nil {
			t.Error(err)
		}
	}
}

// Test the Ping function of the docker helpers
func TestPing(t *testing.T) {
	pong, err := unleash.Ping()
	if err != nil {
		t.Error(err)
	}
	if pong != "OK" {
		t.Error(errors.New("Pong response is not OK but " + pong))
	}
}

// Test the Version function of the Docker helpers
func TestVersion(t *testing.T) {
	version, err := unleash.Version()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error(errors.New("Version is empty !"))
	}
}
