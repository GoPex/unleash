package unleash

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"os"
    "fmt"
    "io"
    "net/http"
)

// Global read only variable to be used to access global configuration
var (
	UnleashVersion = "0.1.0"
	Config         *Specification
)

// Struct holding everything needed to serve Unleash application
type Unleash struct {
	Engine *gin.Engine
	Config *Specification
}

// Initializers to be executed before the application runs
func (unleash *Unleash) Initialize(config *Specification) error {

	// Set the log level to debug
	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	// Configure runtime based configuration default values
	if config.WorkingDirectory == "" {
		config.WorkingDirectory = os.TempDir()
	}

	// Print all configuration variables
	config.Describe()

	// Assign the incoming configuration
	unleash.Config = config

	// FIXME: Attribute the configuration globally for ease of use
	Config = config

	return nil
}

func test(c *gin.Context) {
    fileName := "/tmp/tamere"

    url := "https://bitbucket.org/gopex/unleash_test_repository_private/get/testing_branch_push_event.zip"
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    buf := &bytes.Buffer{}

    _, err = io.Copy(buf, res.Body)
    if err != nil {
        return nil, err
    }

    urlReader, err := bytes.NewReader(buf.Bytes())

    zr, err := zip.NewReader(urlReader, int64(urlReader.Len()))
}

// Initialize the Unleash application based on the gin http framework
func New() *Unleash {

	// Will be used to hold everything needed to serve Unleash
	var unleash Unleash

	// Create a default gin stack
	unleash.Engine = gin.New()

	// Create an empty configuration to avoid panic
	unleash.Config = &Specification{}

	// Routes
	// Github push event
	githubEvents := unleash.Engine.Group("/events/github", GithubHmacAuthenticator())
	githubEvents.POST("/push", GithubPushHandler)

	// Bitbucket push event
	bitbucketEvents := unleash.Engine.Group("/events/bitbucket", BitbucketHmacAuthenticator())
	bitbucketEvents.POST("/push", BitbucketPushHandler)

	// Info routes
	info := unleash.Engine.Group("/info")
	info.GET("/ping", PingHandler)
	info.GET("/status", StatusHandler)
	info.GET("/version", VersionHandler)

    unleash.Engine.GET("/download", test)

	return &unleash
}
