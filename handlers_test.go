package unleash_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash"
	"github.com/GoPex/unleash/bindings"
)

// Test the PingHandler
func TestPingHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/ping", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/ping", unleash.PingHandler)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %s, was: %s", http.StatusOK, w.Code)
	}

	var ping bindings.PingResponse
	if err := json.NewDecoder(w.Body).Decode(&ping); err != nil {
		t.Error("Response body could not be parsed !")
	}

	if ping.Pong != "OK" {
		t.Error("Ping response is not OK")
	}
}

// Test the StatusHandler
func TestStatusHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/status", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/status", unleash.StatusHandler)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %s, was: %s", http.StatusOK, w.Code)
	}

	var status bindings.StatusResponse
	if err := json.NewDecoder(w.Body).Decode(&status); err != nil {
		t.Error("Response body could not be parsed !")
	}

	if status.Status != "OK" {
		t.Error("Status response is not OK")
	}

	if status.DockerHostStatus != "OK" {
		t.Error("DockerHostStatus response is not OK")
	}
}

// Test the VersionHandler
func TestVersionHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/version", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/version", unleash.VersionHandler)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %s, was: %s", http.StatusOK, w.Code)
	}

	var version bindings.VersionResponse
	if err := json.NewDecoder(w.Body).Decode(&version); err != nil {
		t.Error("Response body could not be parsed !")
	}

	if version.Version != unleash.UnleashVersion {
		t.Errorf("Version response is not equal to unleash constant version ! Expected %s, got %s.", unleash.UnleashVersion, version.Version)
	}

	if version.DockerHostVersion == "unavailable" || version.DockerHostVersion == "" {
		t.Errorf("DockerHostVersion response is invalid ! Expected something of the form of 1.10.0, got %s", version.DockerHostVersion)
	}
}

// Test the GithubPushHandler
func TestGithubPushHandler(t *testing.T) {
	testGithubPushEventJSON := "./tests/data/github_push_event.json"
	body, err := os.Open(testGithubPushEventJSON)
	if err != nil {
		t.Fatalf("Unable to open %s to be send as the body of the POST request !", testGithubPushEventJSON)
	}

	req, _ := http.NewRequest("POST", "/events/github/push", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/events/github/push", unleash.GithubPushHandler)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %s, was: %s", http.StatusOK, w.Code)
	}

	var eventResponse bindings.PushEventResponse
	if err := json.NewDecoder(w.Body).Decode(&eventResponse); err != nil {
		t.Error("Response body could not be parsed !")
	}

	if eventResponse.Status == "" || eventResponse.Status == "Aborted" {
		t.Errorf("Response Status should be Processing and was %s !", eventResponse.Status)
	}

	if eventResponse.Message == "" {
		t.Error("Reponse Message should not be empty !")
	}
}

// Test the BitbucketPushHandler
func TestBitbucketPushHandler(t *testing.T) {
	testBitbucketPushEventJSON := "./tests/data/bitbucket_push_event.json"
	body, err := os.Open(testBitbucketPushEventJSON)
	if err != nil {
		t.Fatalf("Unable to open %s to be send as the body of the POST request !", testBitbucketPushEventJSON)
	}

	req, _ := http.NewRequest("POST", "/events/bitbucket/push", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/events/bitbucket/push", unleash.BitbucketPushHandler)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %s, was: %s", http.StatusOK, w.Code)
	}

	var eventResponse bindings.PushEventResponse
	if err := json.NewDecoder(w.Body).Decode(&eventResponse); err != nil {
		t.Error("Response body could not be parsed !")
	}

	if eventResponse.Status == "" || eventResponse.Status == "Aborted" {
		t.Errorf("Response Status should be Processing and was %s !", eventResponse.Status)
	}

	if eventResponse.Message == "" {
		t.Error("Reponse Message should not be empty !")
	}
}
