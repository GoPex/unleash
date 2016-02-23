package unleash

import (
    log "github.com/Sirupsen/logrus"

    // go client for git
    "github.com/libgit2/git2go"
)

// Git helpers function, clone a repository from the specified url,
// to the specified path and checkout it to the specified branch.
func Clone(repositoryUrl string, destinationPath string, branch string) error {
    log.Debug("Cloning the repository ", repositoryUrl, " to ", destinationPath, " on branch ", branch, " ...")

    // Preparing the clone call
    opts := git.CloneOptions{}

    // Cloning the repository
    repo, err := git.Clone(repositoryUrl, destinationPath, &opts)
    if err != nil {
        log.Error("Not possible to clone ", repositoryUrl, " cause: ", err)
        return err
    }

    log.Debug("Repository cloned to ", repo.Workdir(), " !")

    return nil
}
