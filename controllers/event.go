package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/bindings"
	"github.com/GoPex/unleash/helpers"
	"github.com/GoPex/unleash/jobs"
)

// PostGithub is an handler for the POST /events/github/push route.
// Based on the event received from Github, this will schedule a
// BuildAndPushFromRepository background job.
func PostGithub(c *gin.Context) {
	// Parse incomming JSON from github
	var pushEvent bindings.GithubPushEvent
	if c.BindJSON(&pushEvent) == nil {
		// Get the branch name from the JSON "ref" attribute
		tokens := strings.Split(pushEvent.Ref, "/")
		branch := tokens[len(tokens)-1]

		// Evaluate Github push event variables in the direct archive download URL
		repositoryArchiveURL := helpers.EvaluateURL(pushEvent.Repository.ArchiveURL, branch)

		// Launch the build in background
		go jobs.BuildAndPushFromRepository(repositoryArchiveURL,
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

// PostBitbucket is and handler for the POST /events/bitbucket/push route.
// Based on the event received from Bitbucket, this will schedule a
// BuildAndPushFromRepository background job.
func PostBitbucket(c *gin.Context) {
	// Parse incomming JSON from github
	var pushEvent bindings.BitbucketPushEvent
	if c.BindJSON(&pushEvent) == nil {
		for _, change := range pushEvent.Push.Changes {
			// The bitbucket push event doesn't contain a direct archive download URL, we need to build it ourself
			repositoryArchiveURL := pushEvent.Repository.Links.HTML.Href + "/get/" + change.New.Name + ".tar.gz"

			// Launch the build in background
			go jobs.BuildAndPushFromRepository(repositoryArchiveURL,
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
