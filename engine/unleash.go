package engine

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/controllers"
	"github.com/GoPex/unleash/helpers"
)

// Unleash struct holding everything needed to serve Unleash application
type Unleash struct {
	Engine *gin.Engine
	Config *helpers.Specification
}

// Initialize to be executed before the application runs
func (unleash *Unleash) Initialize(config *helpers.Specification) error {

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
	helpers.Config = config

	return nil
}

// New initialize the Unleash application based on the gin http framework
func New() *Unleash {

	// Will be used to hold everything needed to serve Unleash
	var unleash Unleash

	// Create an empty configuration to avoid panic
	unleash.Config = &helpers.Specification{}

	// Create a default gin stack
	unleash.Engine = gin.Default()

	// Routes
	// Github push event
	githubEvents := unleash.Engine.Group("/events/github", HmacAuthenticator(verifyGithubSignature))
	githubEvents.POST("/push", controllers.PostGithub)

	// Bitbucket push event
	bitbucketEvents := unleash.Engine.Group("/events/bitbucket", HmacAuthenticator(verifyBitbucketSignature))
	bitbucketEvents.POST("/push", controllers.PostBitbucket)

	// Ping route
	unleash.Engine.GET("/ping", controllers.GetPing)

	// Info routes
	info := unleash.Engine.Group("/info")
	info.GET("/status", controllers.GetStatus)
	info.GET("/version", controllers.GetVersion)

	return &unleash
}
