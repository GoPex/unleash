package unleash_test

import (
    "os"
    "testing"
    "reflect"

    // Unleash package to test
    "github.com/GoPex/unleash"
)

type configTest struct {
    variableName        string
    configVariableName  string
    expected            string
}

var (
    configTests = []configTest {
        {"UNLEASH_PORT", "Port", "3000"},
        {"UNLEASH_WORKING_DIRECTORY", "WorkingDirectory", "/tmp"},
        {"UNLEASH_REGISTRY_URL", "RegistryURL", "localhost:5000"},
        {"UNLEASH_REGISTRY_USERNAME", "RegistryUsername", "username"},
        {"UNLEASH_REGISTRY_PASSWORD", "RegistryPassword", "password"},
        {"UNLEASH_REGISTRY_EMAIL", "RegistryEmail", "username@example.org"},
        {"UNLEASH_API_KEY", "ApiKey", "supersecret"},
    }
)

// Function to enable table driven test for Specification struct
func getConfigValue(spec *unleash.Specification, field string) string {
    r := reflect.ValueOf(spec)
    f := reflect.Indirect(r).FieldByName(field)
    return f.String()
}

// Test the BuildAndPushFromRepository job
func TestParseConfiguration(t *testing.T) {
    // First, set the required environments variables
    for _, tt := range configTests {
        os.Setenv(tt.variableName, tt.expected)
    }

    // Call the parsing function to test
    config, err := unleash.ParseConfiguration()
    if err != nil {
        t.Error(err)
    }

    // Check that the fields of config hab been correctly filled
    for _, tt := range configTests {
        actual := getConfigValue(&config, tt.configVariableName)

        if actual != tt.expected {
            t.Errorf("Config parsed value for %s: expected %s, actual %s !", tt.variableName, tt.expected, actual)
        }
    }
}
