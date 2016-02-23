package unleash_test

import (
    "os"
    "testing"

    // go client for docker
    "github.com/GoPex/dockerclient"

    // Tar utilities for go
    "github.com/Rolinh/targo"

    // Unleash package to test
    "bitbucket.org/gopex/unleash"
)

var (
    dockerClient, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
)

// Test the BuildFromTar function of the docker helpers
func TestBuildFromTar(t *testing.T) {
    defer dockerClient.RemoveImage(testImageRepository + ":fromTar", true)

    id, err := unleash.BuildFromTar(testRepositoryTarPath, testImageRepository + ":fromTar")
    if err != nil {
        t.Error(err)
    }
    if id == "" {
        t.Error("No id was returned !")
    }
}

// Test the BuildFromDirectory function of the docker helpers
func TestBuildFromDirectory(t *testing.T) {
    defer dockerClient.RemoveImage(testImageRepository + ":fromDirectory", true)
    defer os.RemoveAll(testRepositoryExtracted)

    targo.Extract(testRepositoryExtracted, testRepositoryTarPath)
    id, err := unleash.BuildFromDirectory(testRepositoryExtracted, testImageRepository + ":fromDirectory")
    if err != nil {
        t.Error(err)
    }
    if id == "" {
        t.Error("No id was returned !")
    }
}
