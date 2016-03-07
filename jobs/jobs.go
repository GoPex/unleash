package jobs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/GoPex/unleash/helpers"
)

// BuildAndPushFromRepository builds a Docker image from a Dockerfile coming
// from a git repository and push the builded Docker image to the registry
func BuildAndPushFromRepository(repositoryURL string, repositoryFullName string, branch string, commit string) error {
	contextLogger := log.WithFields(log.Fields{
		"repository_url": repositoryURL,
		"repository":     repositoryFullName,
		"branch":         branch,
		"commit":         commit,
	})
	contextLogger.Info("Build started")

	// Generate a temporary and unique directory that will hold the build
	uniqueWorkingDirectory := filepath.Join(helpers.Config.WorkingDirectory, strconv.FormatInt(time.Now().UnixNano(), 10)+"-"+commit)
	defer os.RemoveAll(uniqueWorkingDirectory)

	// Generate the full working directory path for the sake of explicitness
	fullWorkingDirectory := filepath.Join(uniqueWorkingDirectory, repositoryFullName, branch, commit)
	os.MkdirAll(fullWorkingDirectory, 0600)

	// Clone the repository
	contextLogger.Debug("Downloading sources")
	if err := helpers.ExtractRepository(repositoryURL, fullWorkingDirectory); err != nil {
		contextLogger.Error("Error cloning, cause: ", err)
		return err
	}

	// Construct the tag to use when building the image
	imageRepository := strings.ToLower(repositoryFullName)
	if branch != "master" {
		imageRepository = imageRepository + ":" + branch
	} else {
		imageRepository = imageRepository + ":latest"
	}

	// Add a registry enpoint if specified in the configuration
	if helpers.Config.RegistryURL != "" {
		imageRepository = filepath.Join(helpers.Config.RegistryURL, imageRepository)
	}

	// The repository used to tag the image must be in lower case
	contextLogger.Debug("Building Dockerfile", imageRepository)
	if _, err := helpers.BuildFromDirectory(fullWorkingDirectory, imageRepository, contextLogger); err != nil {
		contextLogger.Error("Error building Dockerfile, cause: ", err)
		return err
	}

	// Push the Docker image to the registry
	contextLogger.Debug("Pushing Docker image ", imageRepository)
	if err := helpers.PushImage(imageRepository); err != nil {
		contextLogger.Error("Error pushing Docker image, cause: ", err)
		return err
	}

	contextLogger.Info("Build finished")

	return nil
}
