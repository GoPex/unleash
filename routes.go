package unleash

import (
    "strings"
    "net/http"

    // Minimalist http framework
    "github.com/gin-gonic/gin"

    // Internal bindings
    "github.com/GoPex/unleash/bindings"
)

// Controller for the POST events/github/push route. Based on the event received from Github, this will schedule a BuildAndPushFromRepository background job.
func githubPush (c *gin.Context){
    // Parse incomming JSON from github
    var json bindings.GithubPushEvent
    if c.BindJSON(&json) == nil {
        // Get the branch name from the JSON "ref" attribute
        tokens := strings.Split(json.Ref, "/")
        branch := tokens[len(tokens) -1 ]
        // Launch the build in background
        go BuildAndPushFromRepository(json.Repository.CloneURL, json.Repository.FullName, branch, json.HeadCommit.Id)
        // Render status OK (200)
        c.JSON(http.StatusOK, gin.H{"status": "Processing Github push event for commit " + json.HeadCommit.Id + " on branch " + branch + " of the repository " + json.Repository.FullName + "."})
    } else {
        // Render status BadRequest (400)
        c.JSON(http.StatusBadRequest, gin.H{"status": "JSON binding failed !"})
    }
}

// Controller for the POST events/bitbucket/push route. Based on the event received from Bitbucket, this will schedule a BuildAndPushFromRepository background job.
func bitbucketPush (c *gin.Context){
    // Parse incomming JSON from github
    var json bindings.BitbucketPushEvent
    if c.BindJSON(&json) == nil {
        for _, change := range json.Push.Changes {
            // Launch the build in background
            go BuildAndPushFromRepository(json.Repository.Links.HTML.Href, json.Repository.FullName, change.New.Name, change.New.Target.Hash)
        }
        // Render status OK (200)
        c.JSON(http.StatusOK, gin.H{"status": "Processing Bitbucket push event for commit"})
    } else {
        // Render status BadRequest (400)
        c.JSON(http.StatusBadRequest, gin.H{"status": "JSON binding failed !"})
    }
}
