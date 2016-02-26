package unleash

import (
    "net/http"
    log "github.com/Sirupsen/logrus"

    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    bitbucketoauth "golang.org/x/oauth2/bitbucket"
)

var (
    // Set ClientId and ClientSecret to
    oauthConf = &oauth2.Config{
        ClientID:     "Ty8e5vpNyPRthNCc96",
        ClientSecret: "fSy2cx7eqm5LQbvTxHYTqTWJztEn5wb7",
        Scopes:       []string{"repository", "webhook"},
        Endpoint:     bitbucketoauth.Endpoint,
    }
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString = "thisshouldberandom"

    // Global access code for oath
    bitbucketCode string
)

// /login
func handleBitbucketLogin (c *gin.Context){
    url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
    c.Redirect(http.StatusTemporaryRedirect, url)
}

// /github_oauth_cb. Called by github after authorization is granted
func bitbucketOAuthCallback (c *gin.Context){
    state := c.Query("state")
    if state != oauthStateString {
        log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
        return
    }

    bitbucketCode = c.Query("code")
    c.JSON(http.StatusOK, "Logged in, your webhook on private repository should now work !")
}

func getBitbucketToken() (*oauth2.Token, error) {
    token, err := oauthConf.Exchange(oauth2.NoContext, bitbucketCode)
    if err != nil {
        return nil, err
    }

    return token, nil
}
