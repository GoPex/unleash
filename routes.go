package unleash

import (
    "strings"
    "net/http"

    // Minimalist http framework
    "github.com/gin-gonic/gin"

    // Internal bindings
    "bitbucket.org/gopex/unleash/bindings"
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
        // Render status OK (200) and
        c.JSON(http.StatusOK, gin.H{"status": "Processing Github push event for commit " + json.HeadCommit.Id + " on branch " + branch + " of the repository " + json.Repository.FullName + "."})
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "JSON binding failed !"})
    }
}
