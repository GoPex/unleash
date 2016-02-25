package unleash

import (
    log "github.com/Sirupsen/logrus"

    // Automatic parse of the configuration
    "github.com/kelseyhightower/envconfig"
)

// Log all variables parsed with envconfig
func (specification *Specification) Describe() {
    log.WithFields(log.Fields{
            "WorkingDirectory": specification.WorkingDirectory,
            "RegistryURL": specification.RegistryURL,
            "RegistryUsername": specification.RegistryUsername,
            "Port": specification.Port,
        }).Info("Starting unleash")
}

// Struct to hold the configuration of the application
type Specification struct {
    Port             string
    WorkingDirectory string `envconfig:"working_directory" default:"/tmp"`
    RegistryURL      string `envconfig:"registry_url"`
    RegistryUsername string `envconfig:"registry_username"`
    RegistryPassword string `envconfig:"registry_password"`
    RegistryEmail    string `envconfig:"registry_email"`
    ApiKey           string `envconfig:"api_key"`
}

// Parse the configuration of Unleash based on environment variables
func ParseConfiguration() (Specification, error) {
    // Gather the configuration
    var config Specification
    if err := envconfig.Process("unleash", &config); err != nil {
        return config, err
    }
    return config, nil
}
