package unleash

import (
    "bufio"
    "os"
    log "github.com/Sirupsen/logrus"
    "regexp"
    "strings"

    // Tar utilities for go
    "github.com/Rolinh/targo"

    // go client for docker
    "github.com/GoPex/dockerclient"
)

    //// Construct the tag to use when building the image
    //image_tag := repository_name
    //if branch != ""{
        //image_tag = image_tag + ":" + branch
    //} else {
        //image_tag = image_tag + ":latest"
    //}
    //log.Debug("Image repository and tag used ", image_tag)

func BuildFromDirectory(directoryPath string, imageRepository string) (string, error) {
    log.Info("Building Dockerfile for directory ", directoryPath, " ...")

    // Path to the tar of the repository
    directoryTarPath := directoryPath + ".tar"
    log.Debug("Path to the tar of the directory ", directoryTarPath)

    // Create a tar from the repository
    targo.Create(directoryTarPath, directoryPath + "/")
    defer os.Remove(directoryTarPath)

    // Build the Dockerfile for the created tar containing the repository
    id, err := BuildFromTar(directoryTarPath, imageRepository)
    if err != nil {
        return "", err
    }

    log.Info("Dockerfile for directory ", directoryPath, " has been builded !")

    return id, nil
}

func BuildFromTar(tarPath string, imageRepository string) (string, error) {
    log.Info("Building Dockerfile for tar ", tarPath, " ...")

    // Initialize a Docker client
    docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

    // Open the tar for Docker
    dockerBuildContext, err := os.Open(tarPath)
    defer dockerBuildContext.Close()

    // Build the image configuration send to Docker
    buildImageConfig := &dockerclient.BuildImage{
            Context:        dockerBuildContext,
            RepoName:       imageRepository,
            SuppressOutput: false,
    }

    // Send the build request to Docker
    reader, err := docker.BuildImage(buildImageConfig)
    if err != nil {
        return "", err
    }

    // Will wait for the last lign of the Docker build command to parse the generated id
    buildSuccessful, err := regexp.Compile("(Successfully built )[a-z0-9]*")

    // Capture the build output and wait its end before continuing, get the image id at the end
    var id string
    rd := bufio.NewScanner(reader)
    for rd.Scan() {
        message := rd.Text()
        log.Debug(message)

        if buildSuccessful.MatchString(message) {
            id = extractId(message)
        }
    }

    // The Docker build image output has ended, check for errors
    if err = rd.Err(); err != nil {
        return "", err
    }

    log.Info("Dockerfile for tar ", tarPath, " has been builded ", imageRepository, " with id ", id," !")

    return id, nil
}

func extractId(message string) string {
    tokens := strings.Split(message, " ")
    lastToken := tokens[len(tokens) - 1]
    return strings.TrimSuffix(lastToken, "\\n\"}")
}
