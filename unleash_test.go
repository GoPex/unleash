package unleash_test

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/GoPex/unleash"
)

// Struct to do a table driven test for routes
type expectedRoute struct {
	method  string
	path    string
	handler interface{}
}

// Var used by table driven test for routes
var (
	expectedRoutes = []expectedRoute{
		{"GET", "/info/ping", unleash.PingHandler},
		{"GET", "/info/status", unleash.StatusHandler},
		{"GET", "/info/version", unleash.VersionHandler},
		{"POST", "/events/github/push", unleash.GithubPushHandler},
		{"POST", "/events/bitbucket/push", unleash.BitbucketPushHandler},
	}
)

// Function that return the name of given function in string
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// Test the initialize function of the application
func TestInitialize(t *testing.T) {
	// Create an instance of the application
	unleash := unleash.New()

	// Test the Initialize function
	if err := unleash.Initialize(&unleashConfigTest); err != nil {
		t.Errorf("Cannot initialize the application, cause: %s !", err.Error())
	}
}

// Test the routes definition of the application
func TestRoutes(t *testing.T) {
	// Create an instance of the application
	unleash := unleash.New()

	// Get routes definine by the application
	routesInfo := unleash.Engine.Routes()

	// Test each routes values (method, path, handler)
	found := false
	for _, expected := range expectedRoutes {
		found = false
		for _, route := range routesInfo {
			if expected.method == route.Method && expected.path == route.Path {
				found = true
				if nameOfFunction(expected.handler) != route.Handler {
					t.Errorf("Route handler doest not match for %s %s, expected %s, actual %s", expected.method, expected.path, expected.handler, route.Handler)
				}
			}
		}

		if !found {
			t.Errorf("No route found for %s %s !", expected.method, expected.path)
		}
	}
}
