package unleash_test

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	// Unleash package to test
	"github.com/GoPex/unleash"
)

var (
	testsDirectory   = "./tests"
	workingDirectory = filepath.Join(testsDirectory, "tmp")
	dataDirectory    = filepath.Join(testsDirectory, "data")

	gopexGithubUrl    = "https://github.com/GoPex"
	gopexBitbucketUrl = "https://bitbucket.org/gopex"

	testRepositoryUrl                      = gopexGithubUrl + "/unleash_test_repository.git"
	testRepositoryUrlBitbucket             = gopexBitbucketUrl + "/unleash_test_repository.git"
	testPrivateRepositoryUrl               = gopexBitbucketUrl + "/unleash_test_repository_private.git"
	testDestinationPath                    = filepath.Join(workingDirectory, "unleash_test_repository")
	testRepositoryTarPath                  = filepath.Join(dataDirectory, "unleash_test_repository.tar")
	testUnknowInstructionRepositoryTarPath = filepath.Join(dataDirectory, "unleash_test_repository_unknown_instruction.tar")
	testNonZeroCodeRepositoryTarPath       = filepath.Join(dataDirectory, "unleash_test_repository_non-zero_code.tar")
	testRepositoryExtracted                = filepath.Join(workingDirectory, "unleash_test_repository_extracted")

	testRepositoryFullName         = "GoPex/unleash_test_repository"
	testImageRepository            = "gopex/unleash_test_repository"
	testRepositoryDefaultBranch    = "master"
	testRepositoryNotDefaultBranch = "testing_branch_push_event"
	testRepositoryCommitId         = "bb9a1688dec2d9d8cb24136a41e9bc62ad1d9675"

    testGithubPushEventJSON = filepath.Join(dataDirectory, "github_push_event.json")
    testBitbucketPushEventJSON = filepath.Join(dataDirectory, "bitbucket_push_event.json")

	contextLogger = log.WithFields(log.Fields{
		"environment": "test",
	})
)

func init() {
    // Force gin in test mode
    gin.SetMode(gin.TestMode)

    // Force logrus to log only from warnings
	log.SetLevel(log.WarnLevel)

    // Mock Unleash configuration
	unleash.Config = &unleash.Specification{WorkingDirectory: workingDirectory,
		RegistryURL:      "localhost:5000",
		RegistryUsername: "albinos",
		RegistryPassword: os.Getenv("UNLEASH_REGISTRY_PASSWORD"),
		RegistryEmail:    "gilles.albin@gmail.com",
		ApiKey:           "supersecret",
		GitUsername:      "albinos",
		GitPassword:      os.Getenv("UNLEASH_GIT_PASSWORD")}
}
