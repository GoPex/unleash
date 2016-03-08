package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/bindings"
	"github.com/GoPex/unleash/controllers"
	"github.com/GoPex/unleash/helpers"
)

// Test the GetPing handler
func TestGetPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/ping", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/ping", controllers.GetPing)

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

// Test the GetStatus handler
func TestGetStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/status", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/status", controllers.GetStatus)

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

// Test the GetVersion handler
func TestGetVersion(t *testing.T) {
	req, _ := http.NewRequest("GET", "/info/version", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/info/version", controllers.GetVersion)

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
