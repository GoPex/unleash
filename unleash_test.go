package unleash_test

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

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

	testGithubPushEventJSON    = filepath.Join(dataDirectory, "github_push_event.json")
	testBitbucketPushEventJSON = filepath.Join(dataDirectory, "bitbucket_push_event.json")

	testDockerRegistryUrl = "localhost:5000"

	contextLogger = log.WithFields(log.Fields{
		"environment": "test",
	})

	unleashConfigTest unleash.Specification
)

func init() {
	// Force gin in test mode
	gin.SetMode(gin.TestMode)

	// Force logrus to log only from warnings
	log.SetLevel(log.WarnLevel)

	// Mock Unleash configuration
	unleashConfigTest = unleash.Specification{WorkingDirectory: workingDirectory,
		RegistryURL:      testDockerRegistryUrl,
		RegistryUsername: "gopextest",
		RegistryPassword: os.Getenv("UNLEASH_REGISTRY_PASSWORD"),
		RegistryEmail:    "gilles.albin@gmail.com",
		ApiKey:           "supersecret",
		GitUsername:      "gopextest",
		GitPassword:      os.Getenv("UNLEASH_GIT_PASSWORD"),
		LogLevel:         "warning"}
	unleash.Config = &unleashConfigTest
}

type expectedRoute struct {
	method  string
	path    string
	handler interface{}
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

var (
	expectedRoutes = []expectedRoute{
		{"GET", "/info/ping", unleash.PingHandler},
		{"GET", "/info/status", unleash.StatusHandler},
		{"GET", "/info/version", unleash.VersionHandler},
		{"POST", "/events/github/push", unleash.GithubPushHandler},
		{"POST", "/events/bitbucket/push", unleash.BitbucketPushHandler},
	}
)

func TestInitialize(t *testing.T) {
	// Create an instance of the application
	unleash := unleash.New()

	// Test the Initialize function
	if err := unleash.Initialize(&unleashConfigTest); err != nil {
		t.Errorf("Cannot initialize the application, cause: %s !", err.Error())
	}
}

func TestRoutes(t *testing.T) {
	// Create an instance of the application
	unleash := unleash.New()

	// Get routes definine by the application
	routesInfo := unleash.Engine.Routes()

	// Test each routes values (method, path, handler)
	found := false
	for _, expected := range expectedRoutes {
		found = false
		for _, route := range routesInfo {
			if expected.method == route.Method && expected.path == route.Path {
				found = true
				if nameOfFunction(expected.handler) != route.Handler {
					t.Errorf("Route handler doest not match for %s %s, expected %s, actual %s", expected.method, expected.path, expected.handler, route.Handler)
				}
			}
		}

		if !found {
			t.Errorf("No route found for %s %s !", expected.method, expected.path)
		}
	}
}
