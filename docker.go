package unleash

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"regexp"
	"strings"
	"github.com/Rolinh/targo"
	"github.com/GoPex/dockerclient"
)

var (
	bashColorRegex, _       = regexp.Compile("\x1b\\[[0-9;]*m")
	buildSuccessfulRegex, _ = regexp.Compile("(Successfully built )[a-z0-9]*")
)

type MessageStream struct {
	Error       string `json:"error"`
	ErrorDetail struct {
		Message string `json:"message"`
	} `json:"errorDetail"`
	ID             string `json:"id"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
	Status string `json:"status"`
	Stream string `json:"stream"`
}

// Send a Get _ping request to the docker daemon
func Ping() (string, error) {
	docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	pong, err := docker.Ping()

	if err != nil {
		return "NOK", err
	}

	return pong, nil
}

// Send a Get version request to the docker daemon
func Version() (string, error) {
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	version, err := docker.Version()
	if err != nil {
		return "", err
	}

	return version.Version, nil
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
	targo.Create(directoryTarPath, directoryPath+"/")
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
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Hold the last stream message in a string in order to extract the id at the end of the build
	var lastStreamMessage string

	// Capture the build output and wait its end before continuing, get the image id at the end
	jsonReader := json.NewDecoder(reader)
	var message MessageStream

	for {
		// Decode incoming stream directly in JSON
		if err := jsonReader.Decode(&message); err == io.EOF {
			break
		} else if err != nil {
			return "", errors.New("Error decoding incoming JSON from Docker, cause: " + err.Error())
		}

		// Check for stream feed
		if message.Stream != "" {
			// Save the stream message for id extraction later on
			lastStreamMessage = cleanMessage(message.Stream)
			contextLogger.Debug(lastStreamMessage)
		}

		// Check for status feed
		if message.Status != "" {
			contextLogger.Debug(message.Status, " ", message.ID, " ", message.Progress)
		}

		// Check for errors feed
		if message.ErrorDetail.Message != "" {
			errorDetailMessage := strings.TrimSpace(message.ErrorDetail.Message)
			contextLogger.Debug(errorDetailMessage)
			return "", errors.New(errorDetailMessage)
		}

		// Reset the message stream struct to avoid keeping old messages
		message = MessageStream{}
	}

	// Extract the ID of the image built from the last stream message of the Docker daemon
	var id string
	if lastStreamMessage != "" && buildSuccessfulRegex.MatchString(lastStreamMessage) {
		id = extractId(lastStreamMessage)
	} else {
		return "", errors.New("No id found at the end of the build")
	}

	return id, nil
}

// Extract the id of a response message coming from the Docker daemon
func extractId(message string) string {
	tokens := strings.Split(message, " ")
	return tokens[len(tokens)-1]
}

// Clean unwanted \n in incoming messages, needed for stream feed becasue
// sometime the \n is not the last character, the color code is
func cleanMessage(s string) string {
	return strings.Replace(s, "\n", "", -1)
}
