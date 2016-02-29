package unleash

import (
	log "github.com/Sirupsen/logrus"
	"os"

	// Minimalist http framework
	"github.com/gin-gonic/gin"
)

// Global read only variable to be used to access global configuration
var (
	Config         *Specification
	UnleashVersion = "0.1.0"
)

// Initializers to be executed before the application runs
func initialize(config *Specification) {
	// Set the log level to debug
	log.SetLevel(log.DebugLevel)

	// Configure runtime based configuration default values
	if config.WorkingDirectory == "" {
		config.WorkingDirectory = os.TempDir()
	}

	// Print all configuration variables
	config.Describe()

	// Assign configuration globally
	Config = config
}

// Launch the application based on the gin http framework
func Run() {
	// Parse the configuration
	config, err := ParseConfiguration()
	if err != nil {
		panic("Not able to parse the configuration ! Cause: " + err.Error())
	}

	// Initialize the application
	initialize(&config)

	// Create a default gin stack
	r := gin.Default()

	// Routes
	// Github push event
	githubEvents := r.Group("/events/github", GithubHmacAuthenticator())
	githubEvents.POST("/push", githubPushHandler)

	// Bitbucket push event
	bitbucketEvents := r.Group("/events/bitbucket", BitbucketHmacAuthenticator())
	bitbucketEvents.POST("/push", bitbucketPushHandler)

	// Info routes
	info := r.Group("/info")
	info.GET("/ping", pingHandler)
	info.GET("/status", statusHandler)
	info.GET("/version", versionHandler)

	// Unleash!
	r.Run(":" + config.Port) // listen and serve on port defined by environment variable UNLEASH_PORT
}
