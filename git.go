package unleash

import (
    // go client for git
    "github.com/libgit2/git2go"
)

// Git helpers function, clone a repository from the specified url,
// to the specified path and checkout it to the specified branch.
func Clone(repositoryUrl string, destinationPath string, branch string) (string, error) {
    // Preparing the clone call
    opts := git.CloneOptions{CheckoutBranch: branch}

    // Cloning the repository
    repo, err := git.Clone(repositoryUrl, destinationPath, &opts)
    if err != nil {
        return "", err
    }

    return repo.Workdir(), nil
}
