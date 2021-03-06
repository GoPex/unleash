package helpers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

var (
	// Config make things easier but is bad
	Config *Specification

	// UnleashVersion
	UnleashVersion = "0.1.0"
)

// Describe will log all variables parsed with envconfig
func (specification *Specification) Describe() {
	log.WithFields(log.Fields{
		"Port":                  specification.Port,
		"LogLevel":              specification.LogLevel,
		"WorkingDirectory":      specification.WorkingDirectory,
		"RegistryURL":           specification.RegistryURL,
		"RegistryUsername":      specification.RegistryUsername,
		"GitUsername":           specification.GitUsername,
		"BitbucketRepositories": specification.BitbucketRepositories,
	}).Info("Unleash initialized !")
}

// Specification to hold the configuration of the application
type Specification struct {
	Port                  string `default:"3000"`
	LogLevel              string `envconfig:"log_level" default:"debug"`
	WorkingDirectory      string `envconfig:"working_directory"`
	RegistryURL           string `envconfig:"registry_url"`
	RegistryUsername      string `envconfig:"registry_username"`
	RegistryPassword      string `envconfig:"registry_password"`
	RegistryEmail         string `envconfig:"registry_email"`
	APIKey                string `envconfig:"api_key"`
	GitUsername           string `envconfig:"git_username"`
	GitPassword           string `envconfig:"git_password"`
	BitbucketRepositories string `envconfig:"bitbucket_repositories"`
}

// ParseConfiguration will parse the configuration of Unleash based on environment variables
func ParseConfiguration() (Specification, error) {
	// Gather the configuration
	var config Specification
	if err := envconfig.Process("unleash", &config); err != nil {
		return config, err
	}
	return config, nil
}
