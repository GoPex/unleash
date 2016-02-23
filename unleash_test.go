package unleash_test

import (
    "path/filepath"

    // Unleash package to test
    "bitbucket.org/gopex/unleash"
)

var (
    testsDirectory = "./tests"
    workingDirectory = filepath.Join(testsDirectory, "tmp")
    dataDirectory = filepath.Join(testsDirectory, "data")

    gopexGithubUrl = "https://github.com/GoPex"

    testRepositoryUrl = gopexGithubUrl + "/unleash_test_repository.git"
    testDestinationPath = filepath.Join(workingDirectory, "unleash_test_repository")
    testRepositoryTarPath = filepath.Join(dataDirectory, "unleash_test_repository.tar")
    testRepositoryExtracted = filepath.Join(workingDirectory, "unleash_test_repository_extracted")

    testImageRepository = "gopex/unleash_test_repository"

    Config = Specification{}
)
