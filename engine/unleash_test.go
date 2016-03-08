package engine_test

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/GoPex/unleash/controllers"
	"github.com/GoPex/unleash/engine"
	"github.com/GoPex/unleash/tests"
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
		{"GET", "/ping", controllers.GetPing},
		{"GET", "/info/status", controllers.GetStatus},
		{"GET", "/info/version", controllers.GetVersion},
		{"POST", "/events/github/push", controllers.PostGithub},
		{"POST", "/events/bitbucket/push", controllers.PostBitbucket},
	}
)

// Function that return the name of given function in string
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// Test the initialize function of the application
func TestInitialize(t *testing.T) {
	// Create an instance of the application
	unleash := engine.New()

	// Test the Initialize function
	if err := unleash.Initialize(&tests.UnleashConfigTest); err != nil {
		t.Errorf("Cannot initialize the application, cause: %s !", err.Error())
	}
}

// Test the routes definition of the application
func TestRoutes(t *testing.T) {
	// Create an instance of the application
	unleash := engine.New()

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
