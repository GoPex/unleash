package unleash

import (
    // go client for git
    "github.com/libgit2/git2go"
)

// Git helpers function, clone a repository from the specified url,
// to the specified path and checkout it to the specified branch.
func Clone(repositoryUrl string, destinationPath string, branch string) (string, error) {

    // Authentication callback
    callbacks := git.RemoteCallbacks{
        CredentialsCallback: makeCredentialsCallback(Config.GitUsername, Config.GitPassword),
    }

    // Preparing the clone call
    opts := git.CloneOptions{
        // Clone the required branch
        CheckoutBranch: branch,
        // Remote callbacks are for auth
        FetchOptions: &git.FetchOptions{
            RemoteCallbacks: callbacks,
        },
    }

    // Cloning the repository
    repo, err := git.Clone(repositoryUrl, destinationPath, &opts)
    if err != nil {
        return "", err
    }

    return repo.Workdir(), nil
}

// Callback called by libgit2 if remote repository ask for an authentication
func makeCredentialsCallback(username, password string) git.CredentialsCallback {
    return func(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
        errCode, cred := git.NewCredUserpassPlaintext(username, password)
        return git.ErrorCode(errCode), &cred
    }
}
