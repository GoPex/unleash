package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GoPex/unleash/bindings"
	"github.com/GoPex/unleash/controllers"
	"github.com/GoPex/unleash/tests"
)

type handlerTest struct {
	jsonInputPath string
	handler       gin.HandlerFunc
}

var (
	testGithubPushEventJSON    = filepath.Join(tests.DataDirectory, "github_push_event.json")
	testBitbucketPushEventJSON = filepath.Join(tests.DataDirectory, "bitbucket_push_event.json")

	handlerTests = []handlerTest{
		{testGithubPushEventJSON, controllers.PostGithub},
		{testBitbucketPushEventJSON, controllers.PostBitbucket},
	}
)

// TestPostHandlers of events
func TestPostHandlers(t *testing.T) {
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
