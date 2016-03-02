package unleash

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/gin-gonic/gin"
)

// Global read only variable to be used to access global configuration
var (
	UnleashVersion = "0.1.0"
    Config *Specification
)

// Struct holding everything needed to serve Unleash application
type Unleash struct {
    Engine *gin.Engine
    Config *Specification
}

// Initializers to be executed before the application runs
func (unleash *Unleash) Initialize(config *Specification) error {

	// Set the log level to debug
	log.SetLevel(log.DebugLevel)

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

	return &unleash
}
