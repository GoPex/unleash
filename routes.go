package unleash

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"

	// Minimalist http framework
	"github.com/gin-gonic/gin"

	// Internal bindings
	"github.com/GoPex/unleash/bindings"
)

// Handler for the POST /events/github/push route. Based on the event received from Github, this will schedule a BuildAndPushFromRepository background job.
func githubPushHandler(c *gin.Context) {
	// Parse incomming JSON from github
	var pushEvent bindings.GithubPushEvent
	if c.BindJSON(&pushEvent) == nil {
		// Get the branch name from the JSON "ref" attribute
		tokens := strings.Split(pushEvent.Ref, "/")
		branch := tokens[len(tokens)-1]
		// Launch the build in background
		go BuildAndPushFromRepository(pushEvent.Repository.CloneURL, pushEvent.Repository.FullName, branch, pushEvent.HeadCommit.Id)
		// Render status OK (200)
		c.JSON(http.StatusOK, gin.H{"status": "Processing Github push event for commit " + pushEvent.HeadCommit.Id + " on branch " + branch + " of the repository " + pushEvent.Repository.FullName + "."})
	} else {
		// Render status BadRequest (400)
		c.JSON(http.StatusBadRequest, gin.H{"status": "JSON binding failed !"})
	}
}

// Handler for the POST /events/bitbucket/push route. Based on the event received from Bitbucket, this will schedule a BuildAndPushFromRepository background job.
func bitbucketPushHandler(c *gin.Context) {
	// Parse incomming JSON from github
	var pushEvent bindings.BitbucketPushEvent
	if c.BindJSON(&pushEvent) == nil {
		for _, change := range pushEvent.Push.Changes {
			// Launch the build in background
			go BuildAndPushFromRepository(pushEvent.Repository.Links.HTML.Href, pushEvent.Repository.FullName, change.New.Name, change.New.Target.Hash)
		}
		// Render status OK (200)
		c.JSON(http.StatusOK, gin.H{"status": "Processing Bitbucket push event"})
	} else {
		// Render status BadRequest (400)
		c.JSON(http.StatusBadRequest, gin.H{"status": "JSON binding failed !"})
	}
}

// Handler for the GET /info/ping route. This will respond by a pong JSON message if the server is alive
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pong": "OK"})
}

// Handler for the GET /info/status route. This will respond  by the status of the server and of the docker host in a JSON message.
func statusHandler(c *gin.Context) {
	pong, err := Ping()
	if err != nil {
		log.Error("Error trying to ping Docker host, cause: ", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "OK", "docker_host_status": pong})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "docker_host_status": pong})
}

// Handler for the GET /info/version route. This will respond a JSON message with the version of Docker running in the Docker host.
func versionHandler(c *gin.Context) {
	version, err := Version()
	if err != nil {
		log.Error("Error trying to get the Docker host version, cause: ", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"version": UnleashVersion, "docker_host_version": "unavailable"})
	}
	c.JSON(http.StatusOK, gin.H{"version": UnleashVersion, "docker_host_version": version})
}
