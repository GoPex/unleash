package unleash

import (
    "os"
    log "github.com/Sirupsen/logrus"
    "path/filepath"
)

func BuildAndPushFromRepository(repositoryUrl string, repositoryFullName string, branch string, commit string) error {
    logMessage := "commit " + commit + " on branch " + branch + " of the repository " + repositoryFullName + ", using " + repositoryUrl
    log.Info("Starting build for ", logMessage , " ...")

    // Generate the temporary working directory path
    fullWorkingDirectory := filepath.Join(Config.WorkingDirectory, repositoryFullName, branch, commit)
    defer os.RemoveAll(fullWorkingDirectory)

    // Clone the repository
    if err := Clone(repositoryUrl, fullWorkingDirectory, branch); err != nil {
        log.Error("Error while cloning for", logMessage, " ! Cause: ", err)
        return err
    }

    // Build the Dockerfile of the repository
    if err := BuildFromDirectory(fullWorkingDirectory, branch); err != nil {
        log.Error("Error while building docker file for ", logMessage, " ! Cause: ", err)
        return err
    }

    // Push the Docker image to the registry
    // TODO

    log.Info("Build finished for ", logMessage, " !")
}
