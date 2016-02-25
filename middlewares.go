package unleash

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/hex"
    "errors"
    "net/http"
    "io/ioutil"
    log "github.com/Sirupsen/logrus"

    // Minimalist http framework
    "github.com/gin-gonic/gin"
)

// Test the sha1 signature headers of the incoming request using
// User-Agent to get the api key to use. This is based on the Github
// signature mechanism.
func verifySignature (c *gin.Context) error {
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

    // Get the user agent
    userAgent := c.Request.Header.Get("User-Agent")
    if userAgent == "" {
        return errors.New("Missing User-Agent required for HMAC verification !")
    }

    // As we already read the body buffer, we need to refill it for further use
    buff := bytes.NewBuffer(body)
    c.Request.Body = ioutil.NopCloser(buff)

    // Construct the expected mac with the api key
    mac := hmac.New(sha1.New, []byte(Config.ApiKey))
    mac.Write(body)
    expectedMAC := mac.Sum(nil)
    expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

    // Secure compare the two hmac
    if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
        return errors.New("HMAC verification failed for User-Agent '" + userAgent + "' !")
    }

    // Singatures matches !
    return nil
}

// Gin compatible middleware to check incoming request's signature
func GithubHmacAuthenticator() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authenticate the request
        if err := verifySignature(c); err != nil {
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
