package tests

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/helpers"
)

// Global variables to be used by all tests of this package
var (
	TestsDirectory   = filepath.Join(os.Getenv("GOPATH"), "src", "github.com/GoPex/unleash", "tests")
	WorkingDirectory = filepath.Join(TestsDirectory, "tmp")
	DataDirectory    = filepath.Join(TestsDirectory, "data")

	GopexGithubURL    = "https://api.github.com/repos/GoPex"
	GopexBitbucketURL = "https://bitbucket.org/gopex"

	TestRepositoryName        = "unleash_test_repository"
	TestPrivateRepositoryName = "unleash_test_repository_private"

	TestRepositoryFullName         = "GoPex/unleash_test_repository"
	TestImageRepository            = "gopex/unleash_test_repository"
	TestRepositoryDefaultBranch    = "master"
	TestRepositoryNotDefaultBranch = "testing_branch_push_event"
	TestRepositoryCommitID         = "bb9a1688dec2d9d8cb24136a41e9bc62ad1d9675"

	TestDockerRegistryURL = "localhost:5000"

	ContextLogger = log.WithFields(log.Fields{
		"environment": "test",
	})

	UnleashConfigTest helpers.Specification
)

// Initialze test for the whole package
func init() {
	// Force gin in test mode
	gin.SetMode(gin.TestMode)

	// Force logrus to log only from warnings
	log.SetLevel(log.WarnLevel)

	// Mock Unleash configuration
	UnleashConfigTest = helpers.Specification{WorkingDirectory: WorkingDirectory,
		RegistryURL:      TestDockerRegistryURL,
		RegistryUsername: "gopextest",
		RegistryPassword: os.Getenv("UNLEASH_REGISTRY_PASSWORD"),
		RegistryEmail:    "gilles.albin@gmail.com",
		APIKey:           "supersecret",
		GitUsername:      "albinos",
		GitPassword:      os.Getenv("UNLEASH_GIT_PASSWORD"),
		LogLevel:         "warning"}
	helpers.Config = &UnleashConfigTest
}

// GithubArchiveURL construct a full url to a Github repository archive
func GithubArchiveURL(repositoryName string, branch string) string {
	return GopexGithubURL + "/" + repositoryName + "/tarball/" + branch
}

// BitbucketArchiveURL construct a full url to a Bitbucket repository archive
func BitbucketArchiveURL(repositoryName string, branch string) string {
	return GopexBitbucketURL + "/" + repositoryName + "/get/" + branch + ".tar.gz"
}
