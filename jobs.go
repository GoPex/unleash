package unleash

import (
    "os"
    log "github.com/Sirupsen/logrus"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

// Jobs to build a Docker image from a Dockerfile coming from a repository and push
// it to the registry
func BuildAndPushFromRepository(repositoryUrl string, repositoryFullName string, branch string, commit string) error {
    logMessage := "commit " + commit + " on branch " + branch + " of the repository " + repositoryFullName + ", using " + repositoryUrl
    log.Info("Starting build for ", logMessage , " ...")

    // Generate a temporary and unique directory that will hold the build
    uniqueWorkingDirectory := filepath.Join(Config.WorkingDirectory, strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + commit)
    defer os.RemoveAll(uniqueWorkingDirectory)

    // Generate the full working directory path for the sake of explicitness
    fullWorkingDirectory := filepath.Join(uniqueWorkingDirectory, repositoryFullName, branch, commit)

    // Clone the repository
    if err := Clone(repositoryUrl, fullWorkingDirectory, branch); err != nil {
        log.Error("Error while cloning for", logMessage, " ! Cause: ", err)
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
    if Config.RegistryURL != "" {
        imageRepository = filepath.Join(Config.RegistryURL, imageRepository)
    }

    // The repository used to tag the image must be in lower case
    if _, err := BuildFromDirectory(fullWorkingDirectory, imageRepository); err != nil {
        log.Error("Error while building docker file for ", logMessage, " ! Cause: ", err)
        return err
    }

    // Push the Docker image to the registry
    if err := PushImage(imageRepository); err != nil {
        log.Error("Error while pushing Docker image for ", logMessage, " ! Cause: ", err)
        return err
    }

    log.Info("Build finished for ", logMessage, " !")

    return nil
}
