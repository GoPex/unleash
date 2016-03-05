package unleash

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	// Internal bindings
	"github.com/GoPex/unleash/bindings"
)

// GithubPushHandler is an handler for the POST /events/github/push route.
// Based on the event received from Github, this will schedule a
// BuildAndPushFromRepository background job.
func GithubPushHandler(c *gin.Context) {
	// Parse incomming JSON from github
	var pushEvent bindings.GithubPushEvent
	if c.BindJSON(&pushEvent) == nil {
		// Get the branch name from the JSON "ref" attribute
		tokens := strings.Split(pushEvent.Ref, "/")
		branch := tokens[len(tokens)-1]

		// Evaluate Github push event variables in the direct archive download URL
		repositoryArchiveURL := EvaluateURL(pushEvent.Repository.ArchiveURL, branch)

		// Launch the build in background
		go BuildAndPushFromRepository(repositoryArchiveURL,
			pushEvent.Repository.FullName,
			branch,
			pushEvent.HeadCommit.ID)

		// Render status OK (200)
		c.JSON(http.StatusOK,
			bindings.PushEventResponse{Status: "Processing",
				Message: "Triggered build for Github push event for commit " + pushEvent.HeadCommit.ID + " on branch " + branch + " of the repository " + pushEvent.Repository.FullName + "."},
		)
	} else {
		// Render status BadRequest (400)
		c.JSON(http.StatusBadRequest,
			bindings.PushEventResponse{Status: "Aborted",
				Message: "JSON binding failed !"},
		)
	}
}

// BitbucketPushHandler is and handler for the POST /events/bitbucket/push route.
// Based on the event received from Bitbucket, this will schedule a
// BuildAndPushFromRepository background job.
func BitbucketPushHandler(c *gin.Context) {
	// Parse incomming JSON from github
	var pushEvent bindings.BitbucketPushEvent
	if c.BindJSON(&pushEvent) == nil {
		for _, change := range pushEvent.Push.Changes {
			// The bitbucket push event doesn't contain a direct archive download URL, we need to build it ourself
			repositoryArchiveURL := pushEvent.Repository.Links.HTML.Href + "/get/" + change.New.Name + ".tar.gz"

			// Launch the build in background
			go BuildAndPushFromRepository(repositoryArchiveURL,
				pushEvent.Repository.FullName,
				change.New.Name,
				change.New.Target.Hash)
		}

		// Render status OK (200)
		c.JSON(http.StatusOK,
			bindings.PushEventResponse{Status: "Processing",
				Message: "Triggered build for Bitbucket push event"},
		)
	} else {
		// Render status BadRequest (400)
		c.JSON(http.StatusBadRequest,
			bindings.PushEventResponse{Status: "Aborted",
				Message: "JSON binding failed !"},
		)
	}
}

// PingHandler is the handler for the GET /info/ping route.
// This will respond by a pong JSON message if the server is alive
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, bindings.PingResponse{Pong: "OK"})
}

// StatusHandler is an Handler for the GET /info/status route.
// This will respond  by the status of the server and of the docker host in a
// JSON message.
func StatusHandler(c *gin.Context) {
	pong, err := Ping()
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

// VersionHandler Handler for the GET /info/version route. This will respond a
// JSON message with the version of Docker running in the Docker host.
func VersionHandler(c *gin.Context) {
	version, err := Version()
	if err != nil {
		log.Error("Error trying to get the Docker host version, cause: ", err)
		c.JSON(http.StatusServiceUnavailable,
			bindings.VersionResponse{Version: UnleashVersion,
				DockerHostVersion: "unavailable"},
		)
	}
	c.JSON(http.StatusOK,
		bindings.VersionResponse{Version: UnleashVersion,
			DockerHostVersion: version},
	)
}
