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
	router := gin.Default()

	// Routes
	// Github push event
	githubEvents := router.Group("/events/github", GithubHmacAuthenticator())
	githubEvents.POST("/push", GithubPushHandler)

	// Bitbucket push event
	bitbucketEvents := router.Group("/events/bitbucket", BitbucketHmacAuthenticator())
	bitbucketEvents.POST("/push", BitbucketPushHandler)

	// Info routes
	info := router.Group("/info")
	info.GET("/ping", PingHandler)
	info.GET("/status", StatusHandler)
	info.GET("/version", VersionHandler)

	// Unleash!
	router.Run(":" + config.Port) // listen and serve on port defined by environment variable UNLEASH_PORT
}
