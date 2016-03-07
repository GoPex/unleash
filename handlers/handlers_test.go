package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/bindings"
	"github.com/GoPex/unleash/handlers"
	"github.com/GoPex/unleash/helpers"
	"github.com/GoPex/unleash/tests"
)

// Test the PingHandler
func TestPingHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/ping", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/ping", handlers.PingHandler)

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
	router.GET("/info/status", handlers.StatusHandler)

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
	router.GET("/info/version", handlers.VersionHandler)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %s, was: %s", http.StatusOK, w.Code)
	}

	var version bindings.VersionResponse
	if err := json.NewDecoder(w.Body).Decode(&version); err != nil {
		t.Error("Response body could not be parsed !")
	}

	if version.Version != helpers.UnleashVersion {
		t.Errorf("Version response is not equal to unleash constant version ! Expected %s, got %s.", helpers.UnleashVersion, version.Version)
	}

	if version.DockerHostVersion == "unavailable" || version.DockerHostVersion == "" {
		t.Errorf("DockerHostVersion response is invalid ! Expected something of the form of 1.10.0, got %s", version.DockerHostVersion)
	}
}

type pushHandler func(c *gin.Context)

type handlerTest struct {
	jsonInputPath string
	handler       gin.HandlerFunc
}

var (
	testGithubPushEventJSON    = filepath.Join(tests.DataDirectory, "github_push_event.json")
	testBitbucketPushEventJSON = filepath.Join(tests.DataDirectory, "bitbucket_push_event.json")

	handlerTests = []handlerTest{
		{testGithubPushEventJSON, handlers.GithubPushHandler},
		{testBitbucketPushEventJSON, handlers.BitbucketPushHandler},
	}
)

// TestPushHandlers
func TestPushHandlers(t *testing.T) {
	for _, handlerTest := range handlerTests {
		body, err := os.Open(handlerTest.jsonInputPath)
		if err != nil {
			t.Fatalf("Unable to open %s to be send as the body of the POST request ! Cause: %s", testGithubPushEventJSON, err)
		}

		req, _ := http.NewRequest("POST", "/push", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.New()
		router.POST("/push", handlerTest.handler)

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
}
