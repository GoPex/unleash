package bindings

// PingResponse is a JSON struct to use when responding to the client
type PingResponse struct {
	Pong string `json:"pong"`
}

// StatusResponse is a JSON struct to use when responding to the client
type StatusResponse struct {
	DockerHostStatus string `json:"docker_host_status"`
	Status           string `json:"status"`
}

// VersionResponse is a JSON struct to use when responding to the client
type VersionResponse struct {
	DockerHostVersion string `json:"docker_host_version"`
	Version           string `json:"version"`
}

// PushEventResponse is a JSON struct to use when responding to the client
type PushEventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
