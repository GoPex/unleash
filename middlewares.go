package unleash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/bindings"
)

// verifySignature is a function to check the incoming's request signature
type verifySignature func(c *gin.Context) error

// HmacAuthenticator is a gin compatible middleware to check incoming request's signature
// for Github
func HmacAuthenticator(verify verifySignature) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authenticate the request
		if err := verify(c); err != nil {
			log.WithFields(log.Fields{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"user-agent": c.Request.UserAgent(),
				"status":     http.StatusUnauthorized,
			}).Warn(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// Test the sha1 signature headers of the incoming request using
// User-Agent to get the api key to use. This is based on the Github
// signature mechanism.
func verifyGithubSignature(c *gin.Context) error {
	// Get the signature
	sig := c.Request.Header.Get("X-Hub-Signature")
	if sig == "" {
		return errors.New("Missing X-Hub-Signature required for HMAC verification !")
	}

	// Read the body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return errors.New("Unable to read the body for HMAC verification !")
	}

	// As we already read the body buffer, we need to refill it for further use
	buff := bytes.NewBuffer(body)
	c.Request.Body = ioutil.NopCloser(buff)

	// Get the user agent
	userAgent := c.Request.Header.Get("User-Agent")
	if userAgent == "" {
		return errors.New("Missing User-Agent required for HMAC verification !")
	}

	// Construct the expected mac with the api key
	mac := hmac.New(sha1.New, []byte(Config.APIKey))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

	// Secure compare the two hmac
	if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
		return errors.New("HMAC verification failed for User-Agent '" + userAgent + "' !")
	}

	// Signatures matches !
	return nil
}

// Test the sha1 signature headers of the incoming request using
// User-Agent to get the api key to use. This is based on the Bitbucket
// signature mechanism.
//
// FIXME: at the moment, we check that the incomming repository is
// one registred in Unleash as Bitbucket don't support body hmac
// signature yet.
func verifyBitbucketSignature(c *gin.Context) error {

	// Read the body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return errors.New("Unable to read the body for HMAC verification !")
	}

	// As we already read the body buffer, we need to refill it for further use
	buff := bytes.NewBuffer(body)
	c.Request.Body = ioutil.NopCloser(buff)

	// Convert payload to JSON
	var pushEvent bindings.BitbucketPushEvent
	err = json.Unmarshal(body, &pushEvent)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return err
	}

	// Check if the repository in the incoming push event is
	// registred
	unauthorized := true
	repositoryToBuild := pushEvent.Repository.Links.HTML.Href
	tokens := strings.Split(Config.BitbucketRepositories, ",")
	for _, registredRepositoryURL := range tokens {
		if repositoryToBuild == registredRepositoryURL {
			unauthorized = false
		}
	}

	if unauthorized {
		return errors.New("Unknown repository url " + repositoryToBuild + " coming from push event !")
	}

	// Signatures matches !
	return nil
}
