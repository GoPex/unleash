package unleash

import (
    "encoding/json"
    "errors"
    "io"
    "os"
    log "github.com/Sirupsen/logrus"
    "regexp"
    "strings"

    // Tar utilities for go
    "github.com/Rolinh/targo"

    // go client for docker
    "github.com/GoPex/dockerclient"
)

var (
    bashColorRegex, _ = regexp.Compile("\x1b\\[[0-9;]*m")
    buildSuccessfulRegex, _ = regexp.Compile("(Successfully built )[a-z0-9]*")
)

type MessageStream struct {
    Stream      string `json:"stream"`
    Error       string `json:"error"`
    ErrorDetail struct {
        Message string `json:"message"`
    } `json:"errorDetail"`
}

// Send a PushImage request to the docker daemon
func PushImage(imageRepository string) error {
    // Initialize a Docker client
    docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

    // Authentication
    authentication := dockerclient.AuthConfig{Username: Config.RegistryUsername, Password: Config.RegistryPassword, Email: Config.RegistryEmail}

    // Push the image to the default registry
    if err := docker.PushImage(imageRepository, "", &authentication); err != nil {
        return err
    }

    return nil
}

// Send a BuildImage request to the docker daemon by sending a tar
// built with the given directory.
func BuildFromDirectory(directoryPath string, imageRepository string, contextLogger *log.Entry) (string, error) {
    // Path to the tar of the repository
    directoryTarPath := directoryPath + ".tar"

    // Create a tar from the repository
    targo.Create(directoryTarPath, directoryPath + "/")
    defer os.Remove(directoryTarPath)

    // Build the Dockerfile for the created tar containing the repository
    id, err := BuildFromTar(directoryTarPath, imageRepository, contextLogger)
    if err != nil {
        return "", err
    }

    return id, nil
}

// Send a BuildImage request to the docker daemon by sending the given tar.
func BuildFromTar(tarPath string, imageRepository string, contextLogger *log.Entry) (string, error) {
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
    defer reader.Close()
    if err != nil {
        return "", err
    }

    // Capture the build output and wait its end before continuing, get the image id at the end
    jsonReader := json.NewDecoder(reader)
    var message MessageStream
    for {
        if err := jsonReader.Decode(&message); err == io.EOF {
            break
        } else if err != nil {
            return "", errors.New("Error decoding incoming JSON from Docker, cause: " + err.Error())
        }

        if message.Stream != "" {
            message.Stream = cleanMessage(message.Stream)
            contextLogger.Debug(message.Stream)
        }
        if message.ErrorDetail.Message != "" {
            message.ErrorDetail.Message = cleanMessage(message.ErrorDetail.Message)
            contextLogger.Debug(message.ErrorDetail.Message)
            return "", errors.New(message.ErrorDetail.Message)
        }
    }

    var id string
    if message.Stream != "" && buildSuccessfulRegex.MatchString(message.Stream) {
        id = extractId(message.Stream)
    }

    return id, nil
}

// Extract the id of a response message coming from the Docker daemon
func extractId(message string) string {
    tokens := strings.Split(message, " ")
    return tokens[len(tokens) - 1]
}

// Clean unwanted characters in incoming messages
func cleanMessage(s string) string {
        safe := bashColorRegex.ReplaceAllString(s, "")
        return strings.TrimSpace(safe)
}
