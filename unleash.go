package unleash

import (
    log "github.com/Sirupsen/logrus"

    // Minimalist http framework
    "github.com/gin-gonic/gin"

    // Automatic parse of the configuration
    "github.com/kelseyhightower/envconfig"
)

var (
    Config Specification
)

type Specification struct {
    WorkingDirectory string `envconfig:"working_directory" default:"/tmp"`
}

func Run() {
    // Set the log level to debug
    log.SetLevel(log.DebugLevel)

    // Gather the configuration
    if err := envconfig.Process("unleash", &Config); err != nil {
        log.Fatal(err)
    }

    // Create a default gin stack
    r := gin.Default()

    // Routes
    r.POST("/events/github/push", githubPush)

    // Unleash!
    r.Run() // listen and serve on port defined by environment variable PORT
}
