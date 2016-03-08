package controllers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/bindings"
	"github.com/GoPex/unleash/helpers"
)

// GetPing is the handler for the GET /info/ping route.
// This will respond by a pong JSON message if the server is alive
func GetPing(c *gin.Context) {
	c.JSON(http.StatusOK, bindings.PingResponse{Pong: "OK"})
}

// GetStatus is an handler for the GET /info/status route.
// This will respond  by the status of the server and of the docker host in a
// JSON message.
func GetStatus(c *gin.Context) {
	pong, err := helpers.Ping()
	if err != nil {
		log.Error("Error trying to ping Docker host, cause: ", err)
		c.JSON(http.StatusServiceUnavailable,
			bindings.StatusResponse{Status: "OK",
				DockerHostStatus: pong},
		)
	}
	c.JSON(http.StatusOK,
		bindings.StatusResponse{Status: "OK",
			DockerHostStatus: pong},
	)
}

// GetVersion is an handler for the GET /info/version route. This will respond a
// JSON message with the version of Docker running in the Docker host.
func GetVersion(c *gin.Context) {
	version, err := helpers.Version()
	if err != nil {
		log.Error("Error trying to get the Docker host version, cause: ", err)
		c.JSON(http.StatusServiceUnavailable,
			bindings.VersionResponse{Version: helpers.UnleashVersion,
				DockerHostVersion: "unavailable"},
		)
	}
	c.JSON(http.StatusOK,
		bindings.VersionResponse{Version: helpers.UnleashVersion,
			DockerHostVersion: version},
	)
}
